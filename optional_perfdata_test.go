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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptionalPerfData_Empty(t *testing.T) {
	// when
	optionalPerfData := OptionalPerfData{}
	actualPerfData, err := optionalPerfData.Get()

	// then
	assert.Error(t, err)
	assert.Empty(t, actualPerfData)
}

func TestNewOptionalPerfData(t *testing.T) {
	// given
	expectedPerfData, _ := NewNumericPerfData("expected", 0, "", nil, nil, nil)

	// when
	optionalPerfData := NewOptionalPerfData(expectedPerfData)
	actualPerfData, err := optionalPerfData.Get()

	// then
	assert.NoError(t, err)
	assert.Equal(t, true, optionalPerfData.Present())
	assert.Equal(t, expectedPerfData, actualPerfData)
}

func TestOptionalPerfData_OrElse(t *testing.T) {
	// given
	expectedPerfData, _ := NewNumericPerfData("expected", 0, "", nil, nil, nil)
	alternativePerfData, _ := NewNumericPerfData("alternative", 0, "", nil, nil, nil)

	// when
	optionalPerfData1 := NewOptionalPerfData(expectedPerfData)
	optionalPerfData2 := OptionalPerfData{}

	// then
	assert.Equal(t, expectedPerfData, optionalPerfData1.OrElse(alternativePerfData))
	assert.Equal(t, alternativePerfData, optionalPerfData2.OrElse(alternativePerfData))
}

func TestOptionalPerfData_Set(t *testing.T) {
	// given
	perfData1, _ := NewNumericPerfData("perfdata 1", 0, "", nil, nil, nil)
	perfData2, _ := NewNumericPerfData("perfdata 2", 0, "", nil, nil, nil)
	alternativePerfData, _ := NewNumericPerfData("alternative", 0, "", nil, nil, nil)

	// when
	optionalPerfData := NewOptionalPerfData(perfData1)
	optionalPerfData.Set(perfData2)

	// then
	assert.Equal(t, true, optionalPerfData.Present())
	assert.Equal(t, perfData2, optionalPerfData.OrElse(alternativePerfData))
}

func TestOptionalPerfData_If(t *testing.T) {
	// given
	var actualPerfDataPtr *PerfData = nil
	var expectedPerfData PerfData
	expectedPerfData, _ = NewNumericPerfData("expected", 0, "", nil, nil, nil)

	// when
	optionalPerfData := NewOptionalPerfData(expectedPerfData)
	optionalPerfData.If(func(perfData PerfData) {
		actualPerfDataPtr = &perfData
	})

	// then
	assert.Equal(t, &expectedPerfData, actualPerfDataPtr)
}
