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

func TestOptionalResult_Empty(t *testing.T) {
	// when
	optionalResult := OptionalResult{}
	actualResult, err := optionalResult.Get()

	// then
	assert.Error(t, err)
	assert.Empty(t, actualResult)
}

func TestNewOptionalResult(t *testing.T) {
	// given
	expectedResult := NewResult()

	// when
	optionalResult := NewOptionalResult(expectedResult)
	actualResult, err := optionalResult.Get()

	// then
	assert.NoError(t, err)
	assert.Equal(t, true, optionalResult.Present())
	assert.Equal(t, expectedResult, actualResult)
}

func TestOptionalResult_OrElse(t *testing.T) {
	// given
	expectedResult := NewResult()
	alternativeResult := NewResult()

	// when
	optionalResult1 := NewOptionalResult(expectedResult)
	optionalResult2 := OptionalResult{}

	// then
	assert.Equal(t, expectedResult, optionalResult1.OrElse(alternativeResult))
	assert.Equal(t, alternativeResult, optionalResult2.OrElse(alternativeResult))
}

func TestOptionalResult_Set(t *testing.T) {
	// given
	result1 := NewResult()
	result2 := NewResult()
	alternativeResult := NewResult()

	// when
	optionalResult := NewOptionalResult(result1)
	optionalResult.Set(result2)

	// then
	assert.Equal(t, true, optionalResult.Present())
	assert.Equal(t, result2, optionalResult.OrElse(alternativeResult))
}

func TestOptionalResult_If(t *testing.T) {
	// given
	var actualResultPtr *Result
	var expectedResult Result
	expectedResult = NewResult()

	// when
	optionalResult := NewOptionalResult(expectedResult)
	optionalResult.If(func(result Result) {
		actualResultPtr = &result
	})

	// then
	assert.Equal(t, &expectedResult, actualResultPtr)
}
