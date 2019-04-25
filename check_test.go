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

func TestBaseCheck_GetSetMeta(t *testing.T) {
	// given
	check := NewCheck("check", NewSummarizer())
	value1 := float64(13.37)
	value2 := "Hello World"

	// when
	check.SetMeta("test-1", value1)
	check.SetMeta("test-2", value2)

	// then
	assert.Equal(t, value1, check.GetMeta("test-1", nil))
	assert.Equal(t, value2, check.GetMeta("test-2", nil))
	assert.Nil(t, check.GetMeta("missing", nil))
}

func TestBaseCheck_AttachContexts(t *testing.T) {
	// given
	context1 := NewStringInfoContext("context 1")
	context2 := NewStringInfoContext("context 2")
	check := NewCheck("check", NewSummarizer())

	// when
	check.AttachContexts(context1)
	check.AttachContexts(context1, context2)

	// then
	assert.Equal(t, 2, len(check.Contexts()))
	assert.Contains(t, check.Contexts(), context1)
	assert.Contains(t, check.Contexts(), context2)
}

func TestBaseCheck_AttachResources(t *testing.T) {
	// given
	resource1 := NewResource()
	resource2 := NewResource()
	check := NewCheck("check", NewSummarizer())

	// when
	check.AttachResources(resource1)
	check.AttachResources(resource1, resource2)

	// then
	assert.Equal(t, 2, len(check.Resources()))
	assert.Contains(t, check.Resources(), resource1)
	assert.Contains(t, check.Resources(), resource2)
}

func TestBaseCheck_State(t *testing.T) {
	// given
	check1 := NewCheck("check 1", NewSummarizer())
	check2 := NewCheck("check 2", NewSummarizer())

	// when
	check1.Results().Add(
		NewResult(ResultState(StateOk())),
		NewResult(ResultState(StateWarning())),
		NewResult(ResultState(StateCritical())),
	)

	// then
	assert.Equal(t, StateCritical(), check1.State())
	assert.Equal(t, StateUnknown(), check2.State())
}

func TestBaseCheck_Summary(t *testing.T) {
	// given
	summarizer := NewSummarizer()
	check1 := NewCheck("check OK", summarizer)
	check2 := NewCheck("check PROBLEM", summarizer)
	check3 := NewCheck("check EMPTY", summarizer)

	check1.Results().Add(NewResult(ResultState(StateOk()), ResultHint("OK")))
	check2.Results().Add(NewResult(ResultState(StateWarning()), ResultHint("WARNING")))

	// when
	summary1 := check1.Summary()
	summary2 := check2.Summary()
	summary3 := check3.Summary()

	// then
	assert.Equal(t, summarizer.Ok(check1), summary1)
	assert.Equal(t, summarizer.Problem(check2), summary2)
	assert.Equal(t, summarizer.Empty(), summary3)
}

func TestBaseCheck_VerboseSummary(t *testing.T) {
	// given
	summarizer := NewSummarizer()
	check := NewCheck("check", summarizer)
	check.Results().Add(
		NewResult(ResultState(StateOk()), ResultHint("This is fine.")),
		NewResult(ResultState(StateWarning()), ResultHint("Works on my machine.")),
		NewResult(ResultHint("Purely informational.")),
	)

	// when
	verboseSummary := check.VerboseSummary()

	// then
	assert.Equal(t, summarizer.Verbose(check), verboseSummary)
}
