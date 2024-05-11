package server

import "github.com/go-playground/validator/v10"

var (
	v *validator.Validate
)

func GetValidator() *validator.Validate {
	if v == nil {
		v = validator.New()
	}
	return v
}

type JSONResponse struct {
	Message any `json:"message,omitempty"`
	Data    any `json:"data,omitempty"`
	Error   any `json:"error"`
}
