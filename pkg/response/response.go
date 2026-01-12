package response

import (
	"encoding/json"
	"net/http"

	"github.com/zchelalo/expense-control-back/pkg/meta"
)

type Envelope[T any] struct {
	Data      T         `json:"data,omitempty"`
	Meta      *meta.Meta `json:"meta,omitempty"`
	Error     *APIError `json:"error,omitempty"`
	RequestID string    `json:"request_id,omitempty"`
}

type APIError struct {
	Code 	  string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func WriteJSON[T any](w http.ResponseWriter, status int, data T, meta *meta.Meta, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope[T]{
		Data:      data,
		Meta:      meta,
		RequestID: requestID,
	})
}

func WriteError(w http.ResponseWriter, status int, err APIError, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope[any]{
		Error:     &err,
		RequestID: requestID,
	})
}
