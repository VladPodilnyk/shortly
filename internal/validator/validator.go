package validator

import (
	"fmt"
	"net/http"

	"short.io/internal/data"
)

type Validator struct {
	Errors       map[string]string
	aliasMaxSize int
}

func New(aliasMaxSize int) *Validator {
	return &Validator{Errors: map[string]string{}, aliasMaxSize: aliasMaxSize}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) ValidateUserRequest(req data.UserRequest) {
	_, err := http.Get(req.Url)
	if err != nil {
		v.addError("url", "broken")
	}

	if len(req.Alias) > v.aliasMaxSize {
		errorMessage := fmt.Sprintf("alias max size exceeded, must be lesser or equal to %d", v.aliasMaxSize)
		v.addError("alias", errorMessage)
	}
}

func (v *Validator) addError(key string, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}
