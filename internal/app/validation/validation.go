package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

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
