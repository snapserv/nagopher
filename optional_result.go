// Code generated by go generate
// This file was generated by robots at 2019-04-22 15:35:59.837356595 +0000 UTC

package nagopher

import (
	"encoding/json"
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

func (o OptionalResult) MarshalJSON() ([]byte, error) {
	if o.Present() {
		return json.Marshal(o.value)
	}
	return json.Marshal(nil)
}

func (o *OptionalResult) UnmarshalJSON(data []byte) error {

	if string(data) == "null" {
		o.value = nil
		return nil
	}

	var value Result

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = &value
	return nil
}