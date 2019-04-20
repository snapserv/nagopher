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

func TestRange_String_Empty(t *testing.T) {
	r, err := ParseRange("")
	assert.Nil(t, err)
	assert.Equal(t, "", r.String())
}

func TestRange_String_ExplicitStartEnd(t *testing.T) {
	r, err := ParseRange("1:9")
	assert.Nil(t, err)
	assert.Equal(t, "1:9", r.String())
}

func TestRange_String_OmitStart(t *testing.T) {
	r, err := ParseRange("9")
	assert.Nil(t, err)
	assert.Equal(t, "9", r.String())
}

func TestRange_String_OmitEnd(t *testing.T) {
	r, err := ParseRange("1:")
	assert.Nil(t, err)
	assert.Equal(t, "1:", r.String())
}

func TestRange_String_NegativeInfinityStart(t *testing.T) {
	r, err := ParseRange("~:10")
	assert.Nil(t, err)
	assert.Equal(t, "~:10", r.String())
}

func TestRange_String_NegativeInfinityEnd(t *testing.T) {
	r, err := ParseRange(":~")
	assert.Nil(t, r)
	assert.NotNil(t, err)
}

func TestRange_String_Invert(t *testing.T) {
	r, err := ParseRange("@1:9")
	assert.Nil(t, err)
	assert.Equal(t, "@1:9", r.String())
}

func TestRange_String_LargeNumberStart(t *testing.T) {
	r, err := ParseRange("4200000000:")
	assert.Nil(t, err)
	assert.Equal(t, "4200000000:", r.String())
}

func TestRange_String_LargeNumberEnd(t *testing.T) {
	r, err := ParseRange("4200000000")
	assert.Nil(t, err)
	assert.Equal(t, "4200000000", r.String())
}

func TestRange_ViolationHint_Normal(t *testing.T) {
	r, err := ParseRange("1:9")
	assert.Nil(t, err)
	assert.Equal(t, "outside range 1:9", r.ViolationHint())
}

func TestRange_ViolationHint_OmitStart(t *testing.T) {
	r, err := ParseRange(":9")
	assert.Nil(t, err)
	assert.Equal(t, "outside range 0:9", r.ViolationHint())
}

func TestRange_ViolationHint_OmitEnd(t *testing.T) {
	r, err := ParseRange("1:")
	assert.Nil(t, err)
	assert.Equal(t, "outside range 1:inf", r.ViolationHint())
}

func TestRange_ViolationHint_NegativeInfinityStart(t *testing.T) {
	r, err := ParseRange("~:1")
	assert.Nil(t, err)
	assert.Equal(t, "outside range -inf:1", r.ViolationHint())
}

func TestRange_ViolationHint_Empty(t *testing.T) {
	r, err := ParseRange("")
	assert.Nil(t, err)
	assert.Equal(t, "outside range 0:inf", r.ViolationHint())
}
