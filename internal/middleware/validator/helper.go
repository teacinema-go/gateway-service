package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	HTTPStatus int
	Message    string
	Fields     map[string]string
}

func (e *ValidationError) Error() string {
	if len(e.Fields) == 0 {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Fields)
}

func DecodeAndValidate[T any](r *http.Request, v *validator.Validate) (T, *ValidationError) {
	var req T

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, &ValidationError{
			HTTPStatus: http.StatusBadRequest,
			Message:    "invalid json body",
			Fields:     nil,
		}
	}

	if err := v.Struct(req); err != nil {
		return req, &ValidationError{
			HTTPStatus: http.StatusUnprocessableEntity,
			Message:    "validation failed",
			Fields:     mapValidationErrors(err, req),
		}
	}

	return req, nil
}

func mapValidationErrors(err error, structValue any) map[string]string {
	fields := make(map[string]string)

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		// Get struct type for reflection
		structType := reflect.TypeOf(structValue)
		if structType.Kind() == reflect.Pointer {
			structType = structType.Elem()
		}

		for _, fe := range validationErrors {
			fields[jsonFieldName(fe, structType)] = formatValidationTag(fe)
		}
	}

	return fields
}

func jsonFieldName(fe validator.FieldError, structType reflect.Type) string {
	fieldName := fe.Field()

	// Make sure we're working with a struct type
	if structType.Kind() != reflect.Struct {
		// Fallback to simple lowercase conversion
		if len(fieldName) > 0 {
			return strings.ToLower(fieldName[:1]) + fieldName[1:]
		}
		return fieldName
	}

	// Get the field by name from a struct type
	field, found := structType.FieldByName(fieldName)
	if !found {
		// Fallback to simple lowercase conversion
		if len(fieldName) > 0 {
			return strings.ToLower(fieldName[:1]) + fieldName[1:]
		}
		return fieldName
	}

	// Get json tag
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" {
		// Parse the JSON tag (format: "field_name, omitempty")
		tagParts := strings.Split(jsonTag, ",")
		if len(tagParts) > 0 && tagParts[0] != "" && tagParts[0] != "-" {
			return tagParts[0]
		}
	}

	// Fallback to simple lowercase conversion
	if len(fieldName) > 0 {
		return strings.ToLower(fieldName[:1]) + fieldName[1:]
	}
	return fieldName
}

func formatValidationTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field is required"
	case "oneof":
		return fmt.Sprintf("must be one of: %s", fe.Param())
	case "numeric":
		return "must contain only numbers"
	case "len":
		return fmt.Sprintf("must be exactly %s characters", fe.Param())
	default:
		return fmt.Sprintf("validation failed: %s", fe.Tag())
	}
}
