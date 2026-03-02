package bus_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hebertzin/cqrs/internal/application/bus"
)

type fakeQuery struct{ ID string }

type fakeQueryHandler struct {
	result interface{}
}

func (h *fakeQueryHandler) Handle(_ context.Context, _ bus.Query) (interface{}, error) {
	return h.result, nil
}

func TestQueryBus_QueryRegisteredHandler(t *testing.T) {
	b := bus.NewQueryBus()
	b.Register("FakeQuery", &fakeQueryHandler{result: "query-result"})

	result, err := b.Query(context.Background(), "FakeQuery", fakeQuery{ID: "1"})

	assert.NoError(t, err)
	assert.Equal(t, "query-result", result)
}

func TestQueryBus_QueryUnregisteredHandler(t *testing.T) {
	b := bus.NewQueryBus()

	_, err := b.Query(context.Background(), "Unknown", fakeQuery{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unknown")
}
