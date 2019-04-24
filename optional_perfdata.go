package nagopher

import (
	"errors"
)

// OptionalPerfData is an optional PerfData.
type OptionalPerfData struct {
	value *PerfData
}

// NewOptionalPerfData creates an optional.OptionalPerfData from a PerfData.
func NewOptionalPerfData(v PerfData) OptionalPerfData {
	return OptionalPerfData{&v}
}

// Set sets the PerfData value.
func (o *OptionalPerfData) Set(v PerfData) {
	o.value = &v
}

// Get returns the PerfData value or an error if not present.
func (o OptionalPerfData) Get() (PerfData, error) {
	if !o.Present() {
		var zero PerfData
		return zero, errors.New("value not present")
	}
	return *o.value, nil
}

// Present returns whether or not the value is present.
func (o OptionalPerfData) Present() bool {
	return o.value != nil
}

// OrElse returns the PerfData value or a default value if the value is not present.
func (o OptionalPerfData) OrElse(v PerfData) PerfData {
	if o.Present() {
		return *o.value
	}
	return v
}

// If calls the function f with the value if the value is present.
func (o OptionalPerfData) If(fn func(PerfData)) {
	if o.Present() {
		fn(*o.value)
	}
}
