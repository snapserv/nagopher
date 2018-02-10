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

type Warning interface {
	Warning() string
}

type warningString struct {
	s string
}

type WarningCollection struct {
	warnings []Warning
}

func NewWarning(text string) Warning {
	return &warningString{text}
}

func (w *warningString) Warning() string {
	return w.s
}

func NewWarningCollection() *WarningCollection {
	return &WarningCollection{}
}

func (c *WarningCollection) Add(warnings ...Warning) {
	c.warnings = append(c.warnings, warnings...)
}

func (c *WarningCollection) GetStrings() []string {
	var results []string
	for _, warning := range c.warnings {
		results = append(results, warning.Warning())
	}

	return results
}
