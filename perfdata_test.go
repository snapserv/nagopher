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
	err, perfData := NewPerfData("i", 42, "", nil, nil, nil)
	assert.Nil(t, err)

	err, output := perfData.BuildOutput()
	assert.Nil(t, err)
	assert.Equal(t, "i=42", output)
}

func TestPerfData_String_QuotedLabel(t *testing.T) {
	err, perfData := NewPerfData("i i", 42, "", nil, nil, nil)
	assert.Nil(t, err)

	err, output := perfData.BuildOutput()
	assert.Nil(t, err)
	assert.Equal(t, "'i i'=42", output)
}

func TestPerfData_String_InvalidLabel_Quotes(t *testing.T) {
	err, _ := NewPerfData("i'", 42, "", nil, nil, nil)
	assert.NotNil(t, err)
}

func TestPerfData_String_InvalidLabel_Equals(t *testing.T) {
	err, _ := NewPerfData("i=", 42, "", nil, nil, nil)
	assert.NotNil(t, err)
}

func TestPerfData_String_ValueUnit(t *testing.T) {
	err, perfData := NewPerfData("i", 42, "%", nil, nil, nil)
	assert.Nil(t, err)

	err, output := perfData.BuildOutput()
	assert.Nil(t, err)
	assert.Equal(t, "i=42%", output)
}

func TestPerfData_String_ValueRange(t *testing.T) {
	err, valueRange := ParseRange("0:100")
	assert.Nil(t, err)
	err, perfData := NewPerfData("i", 42, "", valueRange, nil, nil)
	assert.Nil(t, err)

	err, output := perfData.BuildOutput()
	assert.Nil(t, err)
	assert.Equal(t, "i=42;;;;100", output)
}

func TestPerfData_String_WarningRange(t *testing.T) {
	err, warningRange := ParseRange("-100:")
	assert.Nil(t, err)
	err, perfData := NewPerfData("i", 42, "", nil, warningRange, nil)
	assert.Nil(t, err)

	err, output := perfData.BuildOutput()
	assert.Nil(t, err)
	assert.Equal(t, "i=42;-100:", output)
}

func TestPerfData_String_CriticalRange(t *testing.T) {
	err, criticalRange := ParseRange("@-100:0")
	assert.Nil(t, err)
	err, perfData := NewPerfData("i", 42, "", nil, nil, criticalRange)
	assert.Nil(t, err)

	err, output := perfData.BuildOutput()
	assert.Nil(t, err)
	assert.Equal(t, "i=42;;@-100:0", output)
}
