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
