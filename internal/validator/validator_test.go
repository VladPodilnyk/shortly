package validator

import (
	"fmt"
	"testing"

	"short.io/internal/data"
)

func TestValidator(t *testing.T) {
	validator := New(3)
	validUrl := "https://www.google.com"
	brokeUrl := "https://broken.io"

	req1 := data.UserRequest{Url: validUrl, Alias: ""}
	validator.ValidateUserRequest(req1)
	if len(validator.Errors) > 0 {
		t.Error("unexpected broken request")
	}

	req2 := data.UserRequest{Url: brokeUrl, Alias: "abcd"}
	validator.ValidateUserRequest(req2)
	if len(validator.Errors) != 2 {
		errMsg := ""
		for key, value := range validator.Errors {
			errMsg += fmt.Sprintf("- %s %s\n", key, value)
		}
		t.Errorf("expected two errors, but got:\n%s", errMsg)
	}
}
