package nagopher

import (
	"errors"
)

// OptionalBounds is an optional Bounds.
type OptionalBounds struct {
	value *Bounds
}

// NewOptionalBounds creates an optional.OptionalBounds from a Bounds.
func NewOptionalBounds(v Bounds) OptionalBounds {
	return OptionalBounds{&v}
}

// Set sets the Bounds value.
func (o *OptionalBounds) Set(v Bounds) {
	o.value = &v
}

// Get returns the Bounds value or an error if not present.
func (o OptionalBounds) Get() (Bounds, error) {
	if !o.Present() {
		var zero Bounds
		return zero, errors.New("value not present")
	}
	return *o.value, nil
}

// Present returns whether or not the value is present.
func (o OptionalBounds) Present() bool {
	return o.value != nil
}

// OrElse returns the Bounds value or a default value if the value is not present.
func (o OptionalBounds) OrElse(v Bounds) Bounds {
	if o.Present() {
		return *o.value
	}
	return v
}

// If calls the function f with the value if the value is present.
func (o OptionalBounds) If(fn func(Bounds)) {
	if o.Present() {
		fn(*o.value)
	}
}
