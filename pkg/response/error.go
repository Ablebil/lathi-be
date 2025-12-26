package response

import "github.com/go-playground/validator/v10"

type APIError struct {
	Type    string            `json:"type"`
	Message string            `json:"message"`
	Detail  string            `json:"detail,omitempty"`
	Status  int               `json:"status"`
	Fields  map[string]string `json:"fields,omitempty"` // validation errors
}

func NewAPIError(status int, errType, message, detail string) *APIError {
	return &APIError{
		Type:    errType,
		Message: message,
		Detail:  detail,
		Status:  status,
	}
}

func NewValidationError(err error) *APIError {
	fields := make(map[string]string)
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			fields[fe.Field()] = fe.Tag()
		}
	}
	return &APIError{
		Type:    "validation_error",
		Message: "Validation failed",
		Status:  422,
		Fields:  fields,
	}
}

// error helpers
func ErrInternal(detail string) *APIError {
	return NewAPIError(500, "internal_error", "Internal server error", detail)
}
func ErrNotFound(detail string) *APIError {
	return NewAPIError(404, "not_found", "Resource not found", detail)
}
func ErrUnauthorized(detail string) *APIError {
	return NewAPIError(401, "unauthorized", "Unauthorized", detail)
}
func ErrBadRequest(detail string) *APIError {
	return NewAPIError(400, "bad_request", "Bad request", detail)
}
func ErrConflict(detail string) *APIError {
	return NewAPIError(409, "conflict", "Conflict", detail)
}
