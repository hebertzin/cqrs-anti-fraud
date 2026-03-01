package bus

import (
	"context"
	"fmt"
)

// Command is a marker interface for all commands.
type Command interface{}

// CommandResult is the result returned by a command handler.
type CommandResult interface{}

// CommandHandler processes a specific command type.
type CommandHandler interface {
	Handle(ctx context.Context, command Command) (CommandResult, error)
}

// CommandBus routes commands to their registered handlers.
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
