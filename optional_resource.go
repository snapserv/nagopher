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
