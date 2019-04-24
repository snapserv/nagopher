package nagopher

import (
	"errors"
)

// OptionalMetric is an optional Metric.
type OptionalMetric struct {
	value *Metric
}

// NewOptionalMetric creates an optional.OptionalMetric from a Metric.
func NewOptionalMetric(v Metric) OptionalMetric {
	return OptionalMetric{&v}
}

// Set sets the Metric value.
func (o *OptionalMetric) Set(v Metric) {
	o.value = &v
}

// Get returns the Metric value or an error if not present.
func (o OptionalMetric) Get() (Metric, error) {
	if !o.Present() {
		var zero Metric
		return zero, errors.New("value not present")
	}
	return *o.value, nil
}

// Present returns whether or not the value is present.
func (o OptionalMetric) Present() bool {
	return o.value != nil
}

// OrElse returns the Metric value or a default value if the value is not present.
func (o OptionalMetric) OrElse(v Metric) Metric {
	if o.Present() {
		return *o.value
	}
	return v
}

// If calls the function f with the value if the value is present.
func (o OptionalMetric) If(fn func(Metric)) {
	if o.Present() {
		fn(*o.value)
	}
}
