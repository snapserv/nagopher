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
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWarning(t *testing.T) {
	// when
	formatString := "Value %s %d"
	formatArgs := []interface{}{"is", 13.37}
	warning := NewWarning(formatString, formatArgs...)

	// then
	assert.Equal(t, fmt.Sprintf(formatString, formatArgs...), warning.Warning())
}

func TestWarningCollection_Add(t *testing.T) {
	// given
	warning1 := NewWarning("Hello")
	warning2 := NewWarning("World")
	warnings := NewWarningCollection()

	// when
	warnings.Add(warning1, warning2)

	// then
	assert.Equal(t, 2, len(warnings.Get()))
	assert.Contains(t, warnings.Get(), warning1)
	assert.Contains(t, warnings.Get(), warning2)
}

func TestWarningCollection_GetWarningStrings(t *testing.T) {
	// given
	warning1 := NewWarning("Hello")
	warning2 := NewWarning("World")
	warnings := NewWarningCollection()
	warnings.Add(warning1, warning2)

	// when
	warningStrings := warnings.GetWarningStrings()

	// then
	assert.Equal(t, 2, len(warningStrings))
	assert.Contains(t, warningStrings, warning1.Warning())
	assert.Contains(t, warningStrings, warning2.Warning())
}
