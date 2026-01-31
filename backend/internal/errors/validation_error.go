package errors

type ValidationError struct {
	Err string `json:"error"`
}

func NewValidationError(err string) ValidationError {
	return ValidationError{Err: err}
}

func (v *ValidationError) Error() string {
	return v.Err
}
