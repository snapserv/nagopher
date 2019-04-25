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
