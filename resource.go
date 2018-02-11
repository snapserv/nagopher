/*
 * nagopher - Library for writing Nagios plugins in Go
 * Copyright (C) 2018  Pascal Mathis
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

// Resource represents a interface for all resource types.
type Resource interface {
	Probe(warnings *WarningCollection) ([]Metric, error)
}

// BaseResource represents a generic context from which all other resource types should originate.
type BaseResource struct{}

// NewResource instantiates 'BaseResource'.
func NewResource() *BaseResource {
	return &BaseResource{}
}

// Probe executes all the required check logic and returns one or more 'Metric' objects. Warnings can be passed through
// the provided WarningCollection variable which gets passed. In case any of the check logic fail, the error should be
// returned without panicking. The base resource does not execute any check logic and defaults to returning an empty
// Metric slice without an error.
func (r *BaseResource) Probe(warnings *WarningCollection) ([]Metric, error) {
	return []Metric{}, nil
}
