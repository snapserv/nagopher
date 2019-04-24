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

func TestOptionalState_Empty(t *testing.T) {
	// when
	optionalState := OptionalState{}
	actualState, err := optionalState.Get()

	// then
	assert.Error(t, err)
	assert.Empty(t, actualState)
}

func TestNewOptionalState(t *testing.T) {
	// given
	expectedState := StateOk()

	// when
	optionalState := NewOptionalState(expectedState)
	actualState, err := optionalState.Get()

	// then
	assert.NoError(t, err)
	assert.Equal(t, true, optionalState.Present())
	assert.Equal(t, expectedState, actualState)
}

func TestOptionalState_OrElse(t *testing.T) {
	// given
	expectedState := StateOk()
	alternativeState := StateUnknown()

	// when
	optionalState1 := NewOptionalState(expectedState)
	optionalState2 := OptionalState{}

	// then
	assert.Equal(t, expectedState, optionalState1.OrElse(alternativeState))
	assert.Equal(t, alternativeState, optionalState2.OrElse(alternativeState))
}

func TestOptionalState_Set(t *testing.T) {
	// given
	state1 := StateOk()
	state2 := StateWarning()
	alternativeState := StateUnknown()

	// when
	optionalState := NewOptionalState(state1)
	optionalState.Set(state2)

	// then
	assert.Equal(t, true, optionalState.Present())
	assert.Equal(t, state2, optionalState.OrElse(alternativeState))
}

func TestOptionalState_If(t *testing.T) {
	// given
	var actualStatePtr *State = nil
	var expectedState State
	expectedState = StateOk()

	// when
	optionalState := NewOptionalState(expectedState)
	optionalState.If(func(state State) {
		actualStatePtr = &state
	})

	// then
	assert.Equal(t, &expectedState, actualStatePtr)
}
