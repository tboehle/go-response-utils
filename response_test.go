package response_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboehle/go-response-utils"
	"github.com/tboehle/go-response-utils/internal/testutils"
)

func TestSuccessfulResponse(t *testing.T) {
	const dataKeyName = "data"
	exampleData := testutils.ExampleData{
		Message: "This is a test",
		Count:   42,
	}

	assert := assert.New(t)

	req, err := http.NewRequest(http.MethodGet, "/success", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			// Simulate successful response
			response.With(rw, req, http.StatusOK, exampleData)
		})

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// recorder contains the response body
	body, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Errorf("Error reading the response body: %v", err)
	}
	var resp response.SuccessResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		t.Errorf("Failed to unmarshal the response: %v", err)
	}

	if success := resp.Success; success != true {
		t.Errorf("Handler returned wrong success value: got %v want %v", success, true)
	}

	if payload := resp.Payload; payload == nil {
		t.Errorf("Handler returned a payload of nil: got %v want not nil", nil)
	}

	// Build up expected type
	var recievedData testutils.ExampleData

	err = recievedData.UnmarshalFromInterface(resp.Payload)
	if err != nil {
		t.Errorf("Error while unmarshaling the send data: %v", err)
	}

	if !assert.Equal(exampleData, recievedData) {
		t.Errorf("Data does not match, allthough it should.: got %v want %v", recievedData, exampleData)
	}
}

func TestErrorResponse(t *testing.T) {
	const errorMsg = "Test error"
	req, err := http.NewRequest(http.MethodGet, "/error", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			err := errors.New(errorMsg)
			if err != nil {
				// Simulate error response
				response.WithError(rw, req, http.StatusInternalServerError, err)
				return
			}
		})

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// recorder contains the response body
	body, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Errorf("Error reading the response body: %v", err)
	}
	var resp response.ErrorResponse

	err = json.Unmarshal(body, &resp)
	if err != nil {
		t.Errorf("Failed to unmarshal the response: %v", err)
	}

	if success := resp.Success; success != false {
		t.Errorf("Handler returned wrong success value: got %v want %v", success, false)
	}

	if errStatus := resp.ComponentError.Code; errStatus != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong error code: got %v want %v", errStatus, http.StatusInternalServerError)
	}

	if errMsg := resp.ComponentError.Error; errMsg != errorMsg {
		t.Errorf("Handler returned a error of nil: got %v want %v", errMsg, errorMsg)
	}
}
