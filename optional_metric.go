/*
 * nagopher - Library for writing Nagios plugins in Go
 * Copyright (C) 2018-2019  Pascal Mathis
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
