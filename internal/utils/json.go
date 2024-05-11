package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ParseJsonFromHTTPRequest(r io.ReadCloser, data interface{}) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read data from HTTP request: %w", err)
	}
	if err = r.Close(); err != nil {
		return fmt.Errorf("closing request body: %w", err)
	}
	if err = json.Unmarshal(bytes, data); err != nil {
		return fmt.Errorf("parse json HTTP request: %w", err)
	}
	return nil
}
func WriteJSONHTTPResponse(w http.ResponseWriter, httpCode int, data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("WriteJSONHTTPResponse() Marshal: %w", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)
	_, err = w.Write(j)
	if err != nil {
		return fmt.Errorf("WriteJSONHTTPResponse() Write: %w", err)
	}
	return nil
}
