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

import "fmt"

// Warning represents a single warning cont
type Warning interface {
	Warning() string
}

type baseWarning struct {
	message string
}

// NewWarning instantiates a new warning, passing the format string and an arbitrary amount of arguments to
// fmt.Sprintf() for building the string.
func NewWarning(format string, values ...interface{}) Warning {
	warning := &baseWarning{
		message: fmt.Sprintf(format, values...),
	}

	return warning
}

func (w baseWarning) Warning() string {
	return w.message
}

// WarningCollection collects an arbitrary amount of warnings, which can happen during runtime execution.
type WarningCollection interface {
	Add(warnings ...Warning)
	Get() []Warning
	GetWarningStrings() []string
}

type warningCollection struct {
	warnings []Warning
}

// NewWarningCollection instantiates a new WarningCollection without any items.
func NewWarningCollection() WarningCollection {
	return &warningCollection{}
}

func (wc warningCollection) Add(warnings ...Warning) {
	wc.warnings = append(wc.warnings, warnings...)
}

func (wc warningCollection) Get() []Warning {
	return wc.warnings
}

func (wc warningCollection) GetWarningStrings() []string {
	var results []string

	for _, warning := range wc.warnings {
		results = append(results, warning.Warning())
	}

	return results
}
