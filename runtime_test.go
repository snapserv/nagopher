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
	"strings"
	"testing"
)

type MockResource struct {
	*BaseResource
}

func NewMockResource() *MockResource {
	return &MockResource{
		BaseResource: NewResource(),
	}
}

func (r *MockResource) Probe() []Metric {
	return []Metric{
		NewMetric("usage1", 49.4, "%", nil, "usage"),
		NewMetric("usage2", 92.6, "%", nil, "usage"),
		NewMetric("usage3", 83.1, "%", nil, "usage"),
	}
}

func TestRuntime_Execute_NonVerbose(t *testing.T) {
	check := NewCheck("usage", NewBaseSummary())
	check.AttachResources(NewMockResource())
	check.AttachContexts(
		NewScalarContext("usage", nil, nil),
	)

	expected := CheckResult{
		ExitCode: 0,
		Output:   "USAGE OK - usage1 is 49.4% | usage1=49.4% usage2=92.6% usage3=83.1%",
	}

	result := NewRuntime(false).Execute(check)
	assert.Equal(t, expected, result)
}

func TestRuntime_Execute_Verbose(t *testing.T) {
	check := NewCheck("usage", NewBaseSummary())
	check.AttachResources(NewMockResource())
	check.AttachContexts(
		NewScalarContext("usage", nil, nil),
	)

	expected := CheckResult{
		ExitCode: 0,
		Output: strings.Join([]string{
			"USAGE OK - usage1 is 49.4%",
			" | usage1=49.4%",
			"usage2=92.6%",
			"usage3=83.1%",
		}, "\n"),
	}

	result := NewRuntime(true).Execute(check)
	assert.Equal(t, expected, result)
}
