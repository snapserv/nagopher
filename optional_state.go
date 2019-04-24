package nagopher

import (
	"errors"
)

// OptionalState is an optional State.
type OptionalState struct {
	value *State
}

// NewOptionalState creates an optional.OptionalState from a State.
func NewOptionalState(v State) OptionalState {
	return OptionalState{&v}
}

// Set sets the State value.
func (o *OptionalState) Set(v State) {
	o.value = &v
}

// Get returns the State value or an error if not present.
func (o OptionalState) Get() (State, error) {
	if !o.Present() {
		var zero State
		return zero, errors.New("value not present")
	}
	return *o.value, nil
}

// Present returns whether or not the value is present.
func (o OptionalState) Present() bool {
	return o.value != nil
}

// OrElse returns the State value or a default value if the value is not present.
func (o OptionalState) OrElse(v State) State {
	if o.Present() {
		return *o.value
	}
	return v
}

// If calls the function f with the value if the value is present.
func (o OptionalState) If(fn func(State)) {
	if o.Present() {
		fn(*o.value)
	}
}
