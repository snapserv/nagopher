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

func TestSummary_Ok_ReturnsFirstResult(t *testing.T) {
	results := NewResultCollection()
	results.Add(
		NewResult(StateOk, nil, nil, nil, "Result 1"),
		NewResult(StateOk, nil, nil, nil, "Result 2"),
	)

	summary := NewSummary()
	assert.Equal(t, "Result 1", summary.Ok(results))
}

func TestSummary_Problem_ReturnsMostSignificant(t *testing.T) {
	results := NewResultCollection()
	results.Add(
		NewResult(StateWarning, nil, nil, nil, "Result Warning"),
		NewResult(StateOk, nil, nil, nil, "Result Ok"),
		NewResult(StateCritical, nil, nil, nil, "Result Critical"),
		NewResult(StateOk, nil, nil, nil, "Result Ok"),
	)

	summary := NewSummary()
	assert.Equal(t, "Result Critical", summary.Problem(results))
}

func TestSummary_Verbose(t *testing.T) {
	results := NewResultCollection()
	results.Add(
		NewResult(StateCritical, nil, nil, nil, "Reason 1"),
		NewResult(StateWarning, nil, nil, nil, "Reason 2"),
		NewResult(StateOk, nil, nil, nil, "Must be ignored"),
	)

	summary := NewSummary()
	expected := []string{"critical: Reason 1", "warning: Reason 2"}
	assert.Equal(t, expected, summary.Verbose(results))
}

func TestSummary_Empty(t *testing.T) {
	summary := NewSummary()
	assert.Equal(t, "No check results", summary.Empty())
}
