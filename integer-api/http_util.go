package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[code: %s] %s", e.Code, e.Message)
}

func responseJSON(w http.ResponseWriter, v interface{}) {
	responseJSONWithStatus(w, http.StatusOK, v)
}

func responseJSONWithStatus(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if v == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		fmt.Fprintln(w, err.Error())
	}
}
