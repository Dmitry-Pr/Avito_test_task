// Package errors Description: файл для описания ошибок
package errors

// ErrorResponse - структура ошибки
type ErrorResponse struct {
	Errors string `json:"errors"`
}

// NewErrorResponse - создание ошибки
func NewErrorResponse(errors string) ErrorResponse {
	return ErrorResponse{Errors: errors}
}
