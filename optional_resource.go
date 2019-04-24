package nagopher

import (
	"errors"
)

// OptionalResource is an optional Resource.
type OptionalResource struct {
	value *Resource
}

// NewOptionalResource creates an optional.OptionalResource from a Resource.
func NewOptionalResource(v Resource) OptionalResource {
	return OptionalResource{&v}
}

// Set sets the Resource value.
func (o *OptionalResource) Set(v Resource) {
	o.value = &v
}

// Get returns the Resource value or an error if not present.
func (o OptionalResource) Get() (Resource, error) {
	if !o.Present() {
		var zero Resource
		return zero, errors.New("value not present")
	}
	return *o.value, nil
}

// Present returns whether or not the value is present.
func (o OptionalResource) Present() bool {
	return o.value != nil
}

// OrElse returns the Resource value or a default value if the value is not present.
func (o OptionalResource) OrElse(v Resource) Resource {
	if o.Present() {
		return *o.value
	}
	return v
}

// If calls the function f with the value if the value is present.
func (o OptionalResource) If(fn func(Resource)) {
	if o.Present() {
		fn(*o.value)
	}
}
