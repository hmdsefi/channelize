/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package conn

type fieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func newFieldError(field string, error string) fieldError {
	return fieldError{Field: field, Error: error}
}

// Validation represent validation result that include error message,
// error code, and field validation errors.
type Validation struct {
	Error       string       `json:"error,omitempty"`
	Code        string       `json:"code,omitempty"`
	FieldErrors []fieldError `json:"field_errors,omitempty"`
}

// AddFieldError creates a new fieldError object and adds it to the Validation.FieldErrors
func (v *Validation) AddFieldError(field string, err string) {
	v.FieldErrors = append(v.FieldErrors, newFieldError(field, err))
}

// IsValid returns false if there is error, code, or field errors does exist.
// Otherwise, returns true.
func (v Validation) IsValid() bool {
	return !(v.Error != "" || v.Code != "" || len(v.FieldErrors) != 0)
}
