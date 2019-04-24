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
	"math"
	"testing"
)

func TestNewBounds(t *testing.T) {
	// when
	Bounds := NewBounds(InvertedBounds(true), LowerBound(10), UpperBound(20))
	lowerBound := Bounds.Lower().OrElse(math.NaN())
	upperBound := Bounds.Upper().OrElse(math.NaN())

	// then
	assert.Equal(t, true, Bounds.IsInverted())
	assert.Equal(t, float64(10), lowerBound)
	assert.Equal(t, float64(20), upperBound)
}

func TestNewBounds_Empty(t *testing.T) {
	// when
	Bounds := NewBounds()

	// then
	assert.Equal(t, false, Bounds.IsInverted())
	assert.Equal(t, false, Bounds.Lower().Present())
	assert.Equal(t, false, Bounds.Upper().Present())
}

func TestNewBounds_NagiosRange(t *testing.T) {
	// given
	expectedRanges := map[string]struct {
		inverted   bool
		lowerBound float64
		upperBound float64
	}{
		"":    {false, 0, math.Inf(1)},
		":":   {false, 0, math.Inf(1)},
		"~:":  {false, math.Inf(-1), math.Inf(1)},
		"1:":  {false, 1, math.Inf(1)},
		":1":  {false, 0, 1},
		"1:2": {false, 1, 2},
		"@":   {true, 0, math.Inf(1)},
	}

	// when
	parsedRanges := make(map[string]Bounds)
	for specifier := range expectedRanges {
		Bounds, err := NewBoundsFromNagiosRange(specifier)
		assert.NoError(t, err)

		parsedRanges[specifier] = Bounds
	}

	// then
	for specifier, expectedBounds := range expectedRanges {
		parsedRange, ok := parsedRanges[specifier]

		assert.Equal(t, true, ok)
		assert.Equal(t, expectedBounds.inverted, parsedRange.IsInverted())
		assert.Equal(t, expectedBounds.lowerBound, parsedRange.Lower().OrElse(math.NaN()))
		assert.Equal(t, expectedBounds.upperBound, parsedRange.Upper().OrElse(math.NaN()))
	}
}

func TestNewBounds_NagiosRange_Invalid(t *testing.T) {
	// when
	bounds1, err1 := NewBoundsFromNagiosRange(":~")
	bounds2, err2 := NewBoundsFromNagiosRange("no:float")
	bounds3, err3 := NewBoundsFromNagiosRange("::")

	// then
	assert.Error(t, err1)
	assert.Error(t, err2)
	assert.Error(t, err3)
	assert.Nil(t, bounds1)
	assert.Nil(t, bounds2)
	assert.Nil(t, bounds3)
}

func TestBounds_NagiosRange(t *testing.T) {
	// given
	expectedRanges := map[string]string{
		"":    "",
		":":   "",
		"~":   "~",
		"~:":  "~",
		"@":   "@",
		"@:":  "@",
		"1:":  "1",
		"1:2": "1:2",
		":2":  ":2",
	}

	// when
	describedRanges := make(map[string]string)
	for specifier := range expectedRanges {
		Bounds, err := NewBoundsFromNagiosRange(specifier)
		assert.NoError(t, err)

		describedRanges[specifier] = Bounds.ToNagiosRange()
	}

	// then
	for specifier, expectedRange := range expectedRanges {
		describedRange, ok := describedRanges[specifier]

		assert.Equal(t, true, ok)
		assert.Equal(t, expectedRange, describedRange)
	}
}

func TestBounds_Match(t *testing.T) {
	// given
	Bounds := NewBounds(LowerBound(10), UpperBound(20))

	// when
	trueMatches := []bool{
		Bounds.Match(10),
		Bounds.Match(15),
		Bounds.Match(20),
	}
	falseMatches := []bool{
		Bounds.Match(math.NaN()),
		Bounds.Match(math.Inf(1)),
		Bounds.Match(math.Inf(-1)),
		Bounds.Match(9),
		Bounds.Match(21),
	}

	// then
	assert.Subset(t, []bool{true}, trueMatches)
	assert.Subset(t, []bool{false}, falseMatches)
}

func TestBounds_Match_Inverted(t *testing.T) {
	// given
	Bounds := NewBounds(InvertedBounds(true), LowerBound(10), UpperBound(20))

	// when
	trueMatches := []bool{
		Bounds.Match(9),
		Bounds.Match(21),
	}
	falseMatches := []bool{
		Bounds.Match(math.NaN()),
		Bounds.Match(math.Inf(1)),
		Bounds.Match(math.Inf(-1)),
		Bounds.Match(10),
		Bounds.Match(15),
		Bounds.Match(20),
	}

	// then
	assert.Subset(t, []bool{true}, trueMatches)
	assert.Subset(t, []bool{false}, falseMatches)
}

func TestBounds_StringAndViolationHint(t *testing.T) {
	// when
	bounds1 := NewBounds(LowerBound(10), UpperBound(20))
	bounds2 := NewBounds(LowerBound(10), UpperBound(20), InvertedBounds(true))

	// then
	assert.Equal(t, "inside range 10:20", bounds1.String())
	assert.Equal(t, "outside range 10:20", bounds2.String())
	assert.Equal(t, bounds1.String(), bounds2.ViolationHint())
	assert.Equal(t, bounds2.String(), bounds1.ViolationHint())
}

func TestOptionalBoundsPtr(t *testing.T) {
	// when
	bounds := NewBounds()
	optionalBounds1 := OptionalBounds{}
	optionalBounds2 := NewOptionalBounds(bounds)

	// then
	assert.Nil(t, OptionalBoundsPtr(optionalBounds1))
	assert.Empty(t, &bounds, OptionalBoundsPtr(optionalBounds2))
}
