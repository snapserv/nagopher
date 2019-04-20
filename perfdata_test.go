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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPerfData_String_NormalLabel(t *testing.T) {
	perfData, err := NewPerfData("i", 42, "", nil, nil, nil)
	assert.Nil(t, err)

	output, err := perfData.BuildOutput()
	assert.Nil(t, err)
	assert.Equal(t, "i=42", output)
}

func TestPerfData_String_QuotedLabel(t *testing.T) {
	perfData, err := NewPerfData("i i", 42, "", nil, nil, nil)
	assert.Nil(t, err)

	output, err := perfData.BuildOutput()
	assert.Nil(t, err)
	assert.Equal(t, "'i i'=42", output)
}

func TestPerfData_String_InvalidLabel_Quotes(t *testing.T) {
	perfData, err := NewPerfData("i'", 42, "", nil, nil, nil)
	assert.Nil(t, perfData)
	assert.NotNil(t, err)
}

func TestPerfData_String_InvalidLabel_Equals(t *testing.T) {
	perfData, err := NewPerfData("i=", 42, "", nil, nil, nil)
	assert.Nil(t, perfData)
	assert.NotNil(t, err)
}

func TestPerfData_String_ValueUnit(t *testing.T) {
	perfData, err := NewPerfData("i", 42, "%", nil, nil, nil)
	assert.Nil(t, err)

	output, err := perfData.BuildOutput()
	assert.Nil(t, err)
	assert.Equal(t, "i=42%", output)
}

func TestPerfData_String_ValueRange(t *testing.T) {
	valueRange, err := ParseRange("0:100")
	assert.Nil(t, err)
	perfData, err := NewPerfData("i", 42, "", valueRange, nil, nil)
	assert.Nil(t, err)

	output, err := perfData.BuildOutput()
	assert.Nil(t, err)
	assert.Equal(t, "i=42;;;;100", output)
}

func TestPerfData_String_WarningRange(t *testing.T) {
	warningRange, err := ParseRange("-100:")
	assert.Nil(t, err)
	perfData, err := NewPerfData("i", 42, "", nil, warningRange, nil)
	assert.Nil(t, err)

	output, err := perfData.BuildOutput()
	assert.Nil(t, err)
	assert.Equal(t, "i=42;-100:", output)
}

func TestPerfData_String_CriticalRange(t *testing.T) {
	criticalRange, err := ParseRange("@-100:0")
	assert.Nil(t, err)
	perfData, err := NewPerfData("i", 42, "", nil, nil, criticalRange)
	assert.Nil(t, err)

	output, err := perfData.BuildOutput()
	assert.Nil(t, err)
	assert.Equal(t, "i=42;;@-100:0", output)
}
