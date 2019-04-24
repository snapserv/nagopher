package nagopher

import (
	"errors"
)

// OptionalResult is an optional Result.
type OptionalResult struct {
	value *Result
}

// NewOptionalResult creates an optional.OptionalResult from a Result.
func NewOptionalResult(v Result) OptionalResult {
	return OptionalResult{&v}
}

// Set sets the Result value.
func (o *OptionalResult) Set(v Result) {
	o.value = &v
}

// Get returns the Result value or an error if not present.
func (o OptionalResult) Get() (Result, error) {
	if !o.Present() {
		var zero Result
		return zero, errors.New("value not present")
	}
	return *o.value, nil
}

// Present returns whether or not the value is present.
func (o OptionalResult) Present() bool {
	return o.value != nil
}

// OrElse returns the Result value or a default value if the value is not present.
func (o OptionalResult) OrElse(v Result) Result {
	if o.Present() {
		return *o.value
	}
	return v
}

// If calls the function f with the value if the value is present.
func (o OptionalResult) If(fn func(Result)) {
	if o.Present() {
		fn(*o.value)
	}
}
