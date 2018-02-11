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

type TestRuntimeMockResource struct {
	*BaseResource
}

func NewTestRuntimeMockResource() *TestRuntimeMockResource {
	return &TestRuntimeMockResource{
		BaseResource: NewResource(),
	}
}

func (r *TestRuntimeMockResource) Probe(warnings *WarningCollection) ([]Metric, error) {
	return []Metric{
		NewMetric("usage1", 49.4, "%", nil, "usage"),
		NewMetric("usage2", 92.6, "%", nil, "usage"),
		NewMetric("usage3", 83.1, "|", nil, "usage"),
	}, nil
}

func TestRuntime_Execute_NonVerbose(t *testing.T) {
	check := NewCheck("usage", NewBaseSummary())
	check.AttachResources(NewTestRuntimeMockResource())
	check.AttachContexts(
		NewScalarContext("usage", nil, nil),
	)

	expected := CheckResult{
		ExitCode: 0,
		Output: strings.Join([]string{
			"USAGE OK - usage1 is 49.4% | usage1=49.4% usage2=92.6% usage3=83.1",
			"nagopher: stripped illegal character from string [usage1=49.4% usage2=92.6% usage3=83.1]",
		}, "\n") + "\n",
	}

	result := NewRuntime(false).Execute(check)
	assert.Equal(t, expected, result)
}

func TestRuntime_Execute_Verbose(t *testing.T) {
	warningRange, err := ParseRange("10:80")
	assert.Nil(t, err)

	check := NewCheck("usage", NewBaseSummary())
	check.AttachResources(NewTestRuntimeMockResource())
	check.AttachContexts(
		NewScalarContext("usage", warningRange, nil),
	)

	expected := CheckResult{
		ExitCode: 1,
		Output: strings.Join([]string{
			"USAGE WARNING - usage2 is 92.6% (outside range 10:80)" +
				" | usage1=49.4%;10:80 usage2=92.6%;10:80 usage3=83.1;10:80",
			"warning: usage2 is 92.6% (outside range 10:80)",
			"warning: usage3 is 83.1 (outside range 10:80)",
			"nagopher: stripped illegal character from string [usage1=49.4%;10:80 usage2=92.6%;10:80 usage3=83.1;10:80]",
			"nagopher: stripped illegal character from string [warning: usage3 is 83.1 (outside range 10:80)]",
		}, "\n") + "\n",
	}

	result := NewRuntime(true).Execute(check)
	assert.Equal(t, expected, result)
}
