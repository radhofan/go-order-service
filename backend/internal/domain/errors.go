package domain

import "fmt"

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Detail)
	}
	return e.Message
}

func NewNotFoundError(message, detail string) *AppError {
	return &AppError{
		Code:    404,
		Message: message,
		Detail:  detail,
	}
}

func NewBadRequestError(message, detail string) *AppError {
	return &AppError{
		Code:    400,
		Message: message,
		Detail:  detail,
	}
}

func NewConflictError(message, detail string) *AppError {
	return &AppError{
		Code:    409,
		Message: message,
		Detail:  detail,
	}
}

func NewInternalError(message, detail string) *AppError {
	return &AppError{
		Code:    500,
		Message: message,
		Detail:  detail,
	}
}
