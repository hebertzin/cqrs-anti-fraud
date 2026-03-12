package httpresponse_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hebertzin/cqrs/internal/infrastructure/http/httpresponse"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	httpresponse.JSON(w, http.StatusOK, map[string]string{"key": "value"})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "value", body["key"])
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	httpresponse.Error(w, http.StatusBadRequest, "something went wrong")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body httpresponse.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "Bad Request", body.Error)
	assert.Equal(t, "something went wrong", body.Message)
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	httpresponse.Created(w, map[string]int{"id": 1})
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestOK(t *testing.T) {
	w := httptest.NewRecorder()
	httpresponse.OK(w, "data")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	httpresponse.NoContent(w)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
