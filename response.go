package response

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-http-utils/headers"
	log "github.com/sirupsen/logrus"
)

// SuccessResponse represents the structure of a successful JSON response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Payload interface{} `json:"payload"`
}

// ErrorResponse represents the structure of a error JSON response
type ErrorResponse struct {
	Success        bool           `json:"success,omitempty"`
	ComponentError ComponentError `json:"componentError,omitempty"`
}

// ComponentError represents all the pit error infos if an error appears
type ComponentError struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}

// With returns a response with a payload to a client
func With(w http.ResponseWriter, r *http.Request, status int, payload interface{}) {
	response := SuccessResponse{
		Success: true,
		Payload: payload,
	}
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(response); err != nil {
		log.WithFields(log.Fields{
			"ErrorResponse": response,
			"error":         err,
		}).Error("Error while encoding the SuccessfulResponse")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(headers.ContentType, "application/json")
	w.WriteHeader(status)
	if _, err := io.Copy(w, &buf); err != nil {
		log.WithFields(log.Fields{
			"ErrorResponse": response,
			"error":         err,
		}).Error("Error while copying the SuccessfulResponse to the client")
	}
}

// WithError returns a error response to a client
func WithError(w http.ResponseWriter, r *http.Request, status int, err error) {
	errResponse := ErrorResponse{
		Success: false,
		ComponentError: ComponentError{
			Code:  status,
			Error: err.Error(),
		},
	}
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(errResponse); err != nil {
		log.WithFields(log.Fields{
			"ErrorResponse": errResponse,
			"error":         err,
		}).Error("Error while encoding the ErrorResponse")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(headers.ContentType, "application/json")
	w.WriteHeader(status)

	if _, err := io.Copy(w, &buf); err != nil {
		log.WithFields(log.Fields{
			"ErrorResponse": errResponse,
			"error":         err,
		}).Error("Error while copying the ErrorResponse to the client")
	}
}
