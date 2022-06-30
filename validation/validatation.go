/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package validation

const (
	FieldType     = "type"
	FieldChannels = "channels"
	FieldToken    = "token"
)

type Validator interface {
	Validate() *Result
}

type fieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func newFieldError(field string, error string) fieldError {
	return fieldError{Field: field, Error: error}
}

// Result represent validation result that include error message,
// error code, and field validation errors.
type Result struct {
	Error       string       `json:"error,omitempty"`
	Code        string       `json:"code,omitempty"`
	FieldErrors []fieldError `json:"field_errors,omitempty"`
}

// AddFieldError creates a new fieldError object and adds it to the Result.FieldErrors
func (r *Result) AddFieldError(field string, err string) {
	r.FieldErrors = append(r.FieldErrors, newFieldError(field, err))
}

// IsValid returns false if there is error, code, or field errors does exist.
// Otherwise, returns true.
func (r Result) IsValid() bool {
	return !(r.Error != "" || r.Code != "" || len(r.FieldErrors) != 0)
}
