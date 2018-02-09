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

func TestRange_String_Empty(t *testing.T) {
	assert.Equal(t, "", ParseRange("").String())
}

func TestRange_String_ExplicitStartEnd(t *testing.T) {
	assert.Equal(t, "1:9", ParseRange("1:9").String())
}

func TestRange_String_OmitStart(t *testing.T) {
	assert.Equal(t, "9", ParseRange("9").String())
}

func TestRange_String_OmitEnd(t *testing.T) {
	assert.Equal(t, "1:", ParseRange("1:").String())
}

func TestRange_String_NegativeInfinityStart(t *testing.T) {
	assert.Equal(t, "~:10", ParseRange("~:10").String())
}

func TestRange_String_NegativeInfinityEnd(t *testing.T) {
	// TODO: Implement test case
	// assert.Equal(t, "10:~", ParseRange("10:~").String())
}

func TestRange_String_Invert(t *testing.T) {
	assert.Equal(t, "@1:9", ParseRange("@1:9").String())
}

func TestRange_String_LargeNumberStart(t *testing.T) {
	assert.Equal(t, "4200000000:", ParseRange("4200000000:").String())
}

func TestRange_String_LargeNumberEnd(t *testing.T) {
	assert.Equal(t, "4200000000", ParseRange("4200000000").String())
}

func TestRange_ViolationHint_Normal(t *testing.T) {
	assert.Equal(t, "outside range 1:9", ParseRange("1:9").ViolationHint())
}

func TestRange_ViolationHint_OmitStart(t *testing.T) {
	assert.Equal(t, "outside range 0:9", ParseRange(":9").ViolationHint())
}

func TestRange_ViolationHint_OmitEnd(t *testing.T) {
	assert.Equal(t, "outside range 1:inf", ParseRange("1:").ViolationHint())
}

func TestRange_ViolationHint_NegativeInfinityStart(t *testing.T) {
	assert.Equal(t, "outside range -inf:1", ParseRange("~:1").ViolationHint())
}

func TestRange_ViolationHint_Empty(t *testing.T) {
	assert.Equal(t, "outside range 0:inf", ParseRange("").ViolationHint())
}
