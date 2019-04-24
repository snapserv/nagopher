package nagopher

import (
	"errors"
)

// OptionalContext is an optional Context.
type OptionalContext struct {
	value *Context
}

// NewOptionalContext creates an optional.OptionalContext from a Context.
func NewOptionalContext(v Context) OptionalContext {
	return OptionalContext{&v}
}

// Set sets the Context value.
func (o *OptionalContext) Set(v Context) {
	o.value = &v
}

// Get returns the Context value or an error if not present.
func (o OptionalContext) Get() (Context, error) {
	if !o.Present() {
		var zero Context
		return zero, errors.New("value not present")
	}
	return *o.value, nil
}

// Present returns whether or not the value is present.
func (o OptionalContext) Present() bool {
	return o.value != nil
}

// OrElse returns the Context value or a default value if the value is not present.
func (o OptionalContext) OrElse(v Context) Context {
	if o.Present() {
		return *o.value
	}
	return v
}

// If calls the function f with the value if the value is present.
func (o OptionalContext) If(fn func(Context)) {
	if o.Present() {
		fn(*o.value)
	}
}
