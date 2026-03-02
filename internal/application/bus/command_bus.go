package bus

import (
	"context"
	"fmt"
)

type Command interface{}

type CommandResult interface{}

type CommandHandler interface {
	Handle(ctx context.Context, command Command) (CommandResult, error)
}

type CommandBus struct {
	handlers map[string]CommandHandler
}

func NewCommandBus() *CommandBus {
	return &CommandBus{handlers: make(map[string]CommandHandler)}
}

func (b *CommandBus) Register(commandType string, handler CommandHandler) {
	b.handlers[commandType] = handler
}

func (b *CommandBus) Dispatch(ctx context.Context, commandType string, command Command) (CommandResult, error) {
	handler, ok := b.handlers[commandType]
	if !ok {
		return nil, fmt.Errorf("no handler registered for command: %s", commandType)
	}
	return handler.Handle(ctx, command)
}
