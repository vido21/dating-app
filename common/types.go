package common

import (
	"fmt"
	_ "net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type (
	Route struct {
		Method     string
		Path       string
		Handler    echo.HandlerFunc
		Middleware []echo.MiddlewareFunc
	}

	Controller interface {
		Routes() []Route
	}

	CustomValidator struct {
		Validator *validator.Validate
	}

	ValidationError struct {
		Namespace string
		Field     string
		Tag       string
		Message   string
	}

	ValidationErrors []ValidationError
)

func (ve ValidationErrors) Error() string {
	sErrs := make([]string, len(ve))
	for i, validationError := range ve {
		sErrs[i] = validationError.Message
	}
	return strings.Join(sErrs, "\n")
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errors ValidationErrors
		for _, e := range validationErrors {
			errors = append(errors, ValidationError{
				Namespace: e.Namespace(),
				Field:     e.Field(),
				Tag:       e.Tag(),
				Message:   fmt.Sprintf("Field validation for '%s' failed on the '%s' Tag", e.Field(), e.Tag()),
			})
		}
		return errors
	}
	return nil
}
