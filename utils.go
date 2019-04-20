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

// Warning represents a interface for all warning types.
type Warning interface {
	Warning() string
}

// WarningCollection represents a collection of 0-n warnings.
type WarningCollection struct {
	warnings []Warning
}

// NewWarning instantiates 'Warning' with the given text.
func NewWarning(text string) Warning {
	return &warningString{text}
}

// NewWarningCollection instantiates 'WarningCollection', which by default is empty.
func NewWarningCollection() *WarningCollection {
	return &WarningCollection{}
}

// Add adds one or more 'Warning' objects to the collection.
func (c *WarningCollection) Add(warnings ...Warning) {
	c.warnings = append(c.warnings, warnings...)
}

// GetStrings returns a list containing the string representation of all stored 'Warning' objects.
func (c *WarningCollection) GetStrings() []string {
	var results []string
	for _, warning := range c.warnings {
		results = append(results, warning.Warning())
	}

	return results
}

type warningString struct {
	s string
}

func (w *warningString) Warning() string {
	return w.s
}
