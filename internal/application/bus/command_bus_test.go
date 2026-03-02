package bus_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hebertzin/cqrs/internal/application/bus"
)

type fakeCommand struct{ Value string }

type fakeHandler struct {
	called bool
	err    error
}

func (h *fakeHandler) Handle(_ context.Context, _ bus.Command) (bus.CommandResult, error) {
	h.called = true
	return "result", h.err
}

func TestCommandBus_DispatchRegisteredHandler(t *testing.T) {
	b := bus.NewCommandBus()
	handler := &fakeHandler{}
	b.Register("FakeCommand", handler)

	result, err := b.Dispatch(context.Background(), "FakeCommand", fakeCommand{Value: "test"})

	assert.NoError(t, err)
	assert.Equal(t, "result", result)
	assert.True(t, handler.called)
}

func TestCommandBus_DispatchUnregisteredHandler(t *testing.T) {
	b := bus.NewCommandBus()

	_, err := b.Dispatch(context.Background(), "Unknown", fakeCommand{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unknown")
}

func TestCommandBus_HandlerReturnsError(t *testing.T) {
	b := bus.NewCommandBus()
	expectedErr := errors.New("handler error")
	b.Register("FakeCommand", &fakeHandler{err: expectedErr})

	_, err := b.Dispatch(context.Background(), "FakeCommand", fakeCommand{})

	assert.ErrorIs(t, err, expectedErr)
}
