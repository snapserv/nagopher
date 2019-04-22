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

//go:generate optional -type=Resource
package nagopher

type Resource interface {
	Probe(WarningCollection) ([]Metric, error)
}

type baseResource struct{}

func NewResource() Resource {
	return &baseResource{}
}

func (r baseResource) Probe(warnings WarningCollection) ([]Metric, error) {
	return []Metric{}, nil
}
