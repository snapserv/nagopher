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

func TestResultCollection_Count(t *testing.T) {
	resultCollection := NewResultCollection()
	resultCollection.Add(
		NewResult(StateOk, nil, nil, nil, ""),
		NewResult(StateWarning, nil, nil, nil, ""),
		NewResult(StateCritical, nil, nil, nil, ""),
	)

	assert.Equal(t, 3, resultCollection.Count())
}

func TestResultCollection_MostSignificantState_Normal(t *testing.T) {
	resultCollection := NewResultCollection()
	resultCollection.Add(
		NewResult(StateOk, nil, nil, nil, ""),
		NewResult(StateWarning, nil, nil, nil, ""),
		NewResult(StateCritical, nil, nil, nil, ""),
	)

	assert.Equal(t, StateCritical, resultCollection.MostSignificantState())
}

func TestResultCollection_MostSignificantState_Empty(t *testing.T) {
	resultCollection := NewResultCollection()

	assert.Equal(t, StateUnknown, resultCollection.MostSignificantState())
}
