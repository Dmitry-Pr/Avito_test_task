// Package validation Description: Пакет валидации запросов.
package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator - структура валидатора
type Validator struct {
	validate *validator.Validate
}

// NewValidator - создает новый валидатор
func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// ValidateSendCoinRequest - валидация запроса на отправку монет
func (v *Validator) ValidateSendCoinRequest(req interface{}) error {
	if err := v.validate.Struct(req); err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, err.Error())
		}
		return fmt.Errorf("неправильный формат запроса: %s", strings.Join(errs, ", "))
	}
	return nil
}
