package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/hebertzin/cqrs/internal/application/eventbus"
	"github.com/hebertzin/cqrs/internal/domain/event"
)

const exchangeName = "cqrs.events"

type Factory func([]byte) (event.Event, error)

type envelope struct {
	Type    event.Type      `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventBus struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	mu        sync.RWMutex
	handlers  map[event.Type][]eventbus.Handler
	factories map[event.Type]Factory
	logger    *zap.Logger
}

func NewEventBus(url string, logger *zap.Logger) (*EventBus, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close() //nolint:errcheck
		return nil, fmt.Errorf("rabbitmq open channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		ch.Close()   //nolint:errcheck
		conn.Close() //nolint:errcheck
		return nil, fmt.Errorf("rabbitmq declare exchange: %w", err)
	}

	return &EventBus{
		conn:      conn,
		ch:        ch,
		handlers:  make(map[event.Type][]eventbus.Handler),
		factories: make(map[event.Type]Factory),
		logger:    logger,
	}, nil
}

func (b *EventBus) RegisterFactory(eventType event.Type, factory Factory) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.factories[eventType] = factory
}

func (b *EventBus) Publish(ctx context.Context, e event.Event) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("marshal event payload: %w", err)
	}

	body, err := json.Marshal(envelope{Type: e.GetType(), Payload: payload})
	if err != nil {
		return fmt.Errorf("marshal envelope: %w", err)
	}

	if err := b.ch.PublishWithContext(
		ctx,
		exchangeName,
		string(e.GetType()),
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("rabbitmq publish: %w", err)
	}

	return nil
}

func (b *EventBus) Subscribe(eventType event.Type, handler eventbus.Handler) {
	b.mu.Lock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	b.mu.Unlock()

	queueName := "cqrs." + string(eventType)

	q, err := b.ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		b.logger.Error("rabbitmq declare queue", zap.String("queue", queueName), zap.Error(err))
		return
	}

	if err := b.ch.QueueBind(q.Name, string(eventType), exchangeName, false, nil); err != nil {
		b.logger.Error("rabbitmq bind queue", zap.String("queue", queueName), zap.Error(err))
		return
	}

	msgs, err := b.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		b.logger.Error("rabbitmq consume", zap.String("queue", queueName), zap.Error(err))
		return
	}

	go b.consume(eventType, msgs)
}

func (b *EventBus) consume(eventType event.Type, msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		if err := b.dispatch(eventType, msg); err != nil {
			b.logger.Error("rabbitmq dispatch error",
				zap.String("event_type", string(eventType)),
				zap.Error(err),
			)
			msg.Nack(false, true) //nolint:errcheck
			continue
		}
		msg.Ack(false) //nolint:errcheck
	}
}

func (b *EventBus) dispatch(eventType event.Type, msg amqp.Delivery) error {
	var env envelope
	if err := json.Unmarshal(msg.Body, &env); err != nil {
		return fmt.Errorf("unmarshal envelope: %w", err)
	}

	b.mu.RLock()
	factory, hasFactory := b.factories[eventType]
	handlers := b.handlers[eventType]
	b.mu.RUnlock()

	if !hasFactory {
		return fmt.Errorf("no factory registered for event type %q", eventType)
	}

	e, err := factory(env.Payload)
	if err != nil {
		return fmt.Errorf("deserialize event: %w", err)
	}

	ctx := context.Background()
	for _, h := range handlers {
		if err := h(ctx, e); err != nil {
			b.logger.Error("event handler error",
				zap.String("event_type", string(eventType)),
				zap.Error(err),
			)
		}
	}

	return nil
}

func (b *EventBus) Close() error {
	if err := b.ch.Close(); err != nil {
		return fmt.Errorf("close rabbitmq channel: %w", err)
	}
	if err := b.conn.Close(); err != nil {
		return fmt.Errorf("close rabbitmq connection: %w", err)
	}
	return nil
}
