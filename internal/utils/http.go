package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Meta struct {
	Total int32 `json:"total"`
	Limit int32 `json:"limit"`
	Page  int32 `json:"page"`
}
type JSONResponse struct {
	// Body
	Data interface{} `json:"data,omitempty"`
	// Meta
	Meta *Meta `json:"meta,omitempty"`
	//// Error information
	//Error interface{} `json:"error,omitempty"`
}

func Response(w http.ResponseWriter, respBody interface{}) error {
	resp := JSONResponse{
		Data: respBody,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("can't marshal response: %w", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("error sending response: %w", err)
	}
	return nil
}
func ResponseWithMeta(w http.ResponseWriter, respBody interface{}, meta *Meta) error {
	resp := JSONResponse{
		Data: respBody,
		Meta: meta,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("can't marshal response: %w", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("error sending response: %w", err)
	}
	return nil
}

func GetInt32FromQuery(r *http.Request, key string, defaultValue int32) int32 {
	strValue := r.URL.Query().Get(key)
	value, err := strconv.ParseInt(strValue, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int32(value)
}
