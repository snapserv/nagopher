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
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type mockResource struct {
	Resource
}

type mockEmptyResource struct {
	Resource
}

type mockProbeErrorResource struct {
	Resource
}

type mockPerformanceErrorResource struct {
	Resource
}

func TestBaseRuntime_Execute(t *testing.T) {
	// given
	warningThreshold := NewBounds(LowerBound(10), UpperBound(80))
	check1 := NewCheck("usage", NewSummarizer())
	check2 := NewCheck("usage", NewSummarizer())

	check1.AttachResources(newMockResource())
	check1.AttachContexts(NewScalarContext("usage", nil, nil))
	check2.AttachResources(newMockResource())
	check2.AttachContexts(NewScalarContext("usage", &warningThreshold, nil))

	// when
	result1 := NewRuntime(false).Execute(check1) // non-verbose
	result2 := NewRuntime(true).Execute(check2)  // verbose

	// then
	assert.Equal(t, StateOk().ExitCode(), result1.ExitCode())
	assert.Equal(t, strings.Join([]string{
		"USAGE OK - usage1 is 49.4% | usage1=49.4% usage2=92.6% usage3=83.1",
		"nagopher: stripped illegal character from string [usage3=83.1]",
	}, "\n")+"\n", result1.Output())

	assert.Equal(t, StateWarning().ExitCode(), result2.ExitCode())
	assert.Equal(t, strings.Join([]string{
		"USAGE WARNING - usage2 is 92.6% (outside range 10:80) | usage1=49.4%;10:80 usage2=92.6%;10:80 usage3=83.1;10:80",
		"warning: usage2 is 92.6% (outside range 10:80)",
		"warning: usage3 is 83.1 (outside range 10:80)",
		"nagopher: stripped illegal character from string [usage3=83.1;10:80]",
		"nagopher: stripped illegal character from string [warning: usage3 is 83.1 (outside range 10:80)]",
	}, "\n")+"\n", result2.Output())
}

func TestBaseRuntime_Execute_MissingContext(t *testing.T) {
	// given
	check := NewCheck("check", NewSummarizer())
	check.AttachResources(newMockResource())

	// when
	result := NewRuntime(false).Execute(check)

	// then
	assert.Equal(t, StateUnknown().ExitCode(), result.ExitCode())
	assert.Equal(t, strings.Join([]string{
		"CHECK UNKNOWN - nagopher: missing context with name [usage]",
	}, "\n")+"\n", result.Output())
}

func TestBaseRuntime_Execute_Empty(t *testing.T) {
	// given
	check := NewCheck("check", NewSummarizer())
	check.AttachResources(newMockEmptyResource())

	// when
	result := NewRuntime(false).Execute(check)

	// then
	assert.Equal(t, StateUnknown().ExitCode(), result.ExitCode())
	assert.Equal(t, strings.Join([]string{
		"CHECK UNKNOWN - nagopher: resource [*nagopher.mockEmptyResource] did not return any metrics",
	}, "\n")+"\n", result.Output())
}

func TestBaseRuntime_Execute_ProbeError(t *testing.T) {
	// given
	check := NewCheck("check", NewSummarizer())
	check.AttachResources(newMockProbeErrorResource())

	// when
	result := NewRuntime(false).Execute(check)

	// then
	assert.Equal(t, StateUnknown().ExitCode(), result.ExitCode())
	assert.Equal(t, strings.Join([]string{
		"CHECK UNKNOWN - artificial error happened here!",
	}, "\n")+"\n", result.Output())
}

func TestBaseRuntime_Execute_PerformanceError(t *testing.T) {
	// given
	check := NewCheck("check", NewSummarizer())
	check.AttachResources(newMockPerformanceErrorResource())
	check.AttachContexts(NewScalarContext("usage", nil, nil))

	// when
	result := NewRuntime(false).Execute(check)

	// then
	assert.Equal(t, StateUnknown().ExitCode(), result.ExitCode())
	assert.Equal(t, strings.Join([]string{
		"CHECK UNKNOWN - nagopher: collecting performance data failed with [perfdata metric name [inv'=alid] contains invalid characters]",
	}, "\n")+"\n", result.Output())
}

func TestBaseRuntime_ExecuteAndExit(t *testing.T) {
	var resultExitCode int = -1
	var resultOutput string = ""

	// given
	check := NewCheck("check", NewSummarizer())
	check.AttachResources(newMockResource())
	check.AttachContexts(NewScalarContext("usage", nil, nil))

	resultExitFunction = func(exitCode int) { resultExitCode = exitCode }
	resultOutputFunction = func(values ...interface{}) (int, error) {
		resultOutput = fmt.Sprint(values...)
		return len(resultOutput), nil
	}

	// when
	expectedResult := NewRuntime(false).Execute(check)
	NewRuntime(false).ExecuteAndExit(check)

	// then
	assert.Equal(t, int(expectedResult.ExitCode()), resultExitCode)
	assert.Equal(t, expectedResult.Output(), resultOutput)
}

func TestNewCheckResult(t *testing.T) {
	// when
	exitCode := StateOk().ExitCode()
	description := StateOk().Description()
	checkResult := NewCheckResult(exitCode, description)

	// then
	assert.Equal(t, exitCode, checkResult.ExitCode())
	assert.Equal(t, description, checkResult.Output())
}

func newMockResource() Resource {
	return &mockResource{
		Resource: NewResource(),
	}
}

func (r mockResource) Probe(warnings WarningCollection) ([]Metric, error) {
	return []Metric{
		MustNewNumericMetric("usage1", 49.4, "%", nil, "usage"),
		MustNewNumericMetric("usage2", 92.6, "%", nil, "usage"),
		MustNewNumericMetric("usage3", 83.1, "|", nil, "usage"),
	}, nil
}

func newMockEmptyResource() Resource {
	return &mockEmptyResource{
		Resource: NewResource(),
	}
}

func (r mockEmptyResource) Probe(warnings WarningCollection) ([]Metric, error) {
	return []Metric{}, nil
}

func newMockProbeErrorResource() Resource {
	return &mockProbeErrorResource{
		Resource: NewResource(),
	}
}

func (r mockProbeErrorResource) Probe(warnings WarningCollection) ([]Metric, error) {
	return []Metric{}, fmt.Errorf("artificial error happened here!")
}

func newMockPerformanceErrorResource() Resource {
	return &mockPerformanceErrorResource{
		Resource: NewResource(),
	}
}

func (r mockPerformanceErrorResource) Probe(warnings WarningCollection) ([]Metric, error) {
	return []Metric{
		MustNewNumericMetric("inv'=alid", 49.4, "%", nil, "usage"),
	}, nil
}
