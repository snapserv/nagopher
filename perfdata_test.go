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

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPerfData_String_NormalLabel(t *testing.T) {
	perfData := NewPerfData("i", 42, "", nil, nil, nil)
	assert.Equal(t, "i=42", perfData.String())
}

func TestPerfData_String_QuotedLabel(t *testing.T) {
	perfData := NewPerfData("i i", 42, "", nil, nil, nil)
	assert.Equal(t, "'i i'=42", perfData.String())
}

func TestPerfData_String_InvalidLabel_Quotes(t *testing.T) {
	// TODO: Implement test case
}

func TestPerfData_String_InvalidLabel_Equals(t *testing.T) {
	// TODO: Implement test case
}

func TestPerfData_String_ValueUnit(t *testing.T) {
	perfData := NewPerfData("i", 42, "%", nil, nil, nil)
	assert.Equal(t, "i=42%", perfData.String())
}

func TestPerfData_String_ValueRange(t *testing.T) {
	valueRange := ParseRange("0:100")
	perfData := NewPerfData("i", 42, "", valueRange, nil, nil)

	assert.Equal(t, "i=42;;;;100", perfData.String())
}

func TestPerfData_String_WarningRange(t *testing.T) {
	warningRange := ParseRange("-100:")
	perfData := NewPerfData("i", 42, "", nil, warningRange, nil)

	assert.Equal(t, "i=42;-100:", perfData.String())
}

func TestPerfData_String_CriticalRange(t *testing.T) {
	criticalRange := ParseRange("@-100:0")
	perfData := NewPerfData("i", 42, "", nil, nil, criticalRange)

	assert.Equal(t, "i=42;;@-100:0", perfData.String())
}
