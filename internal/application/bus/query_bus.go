package bus

import (
	"context"
	"fmt"
)

// Query is a marker interface for all queries.
type Query interface{}

// QueryHandler processes a specific query type and returns a result.
type QueryHandler interface {
	Handle(ctx context.Context, query Query) (interface{}, error)
}

// QueryBus routes queries to their registered handlers.
type QueryBus struct {
	handlers map[string]QueryHandler
}

func NewQueryBus() *QueryBus {
	return &QueryBus{handlers: make(map[string]QueryHandler)}
}

func (b *QueryBus) Register(queryType string, handler QueryHandler) {
	b.handlers[queryType] = handler
}

func (b *QueryBus) Query(ctx context.Context, queryType string, query Query) (interface{}, error) {
	handler, ok := b.handlers[queryType]
	if !ok {
		return nil, fmt.Errorf("no handler registered for query: %s", queryType)
	}
	return handler.Handle(ctx, query)
}
