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

func TestBaseContext_Name(t *testing.T) {
	context := NewContext("ctx", "")
	assert.Equal(t, "ctx", context.Name())
}

func TestBaseContext_Describe_Empty(t *testing.T) {
	context := NewContext("ctx", "")
	metric := NewMetric("metric", 0, "", nil, "")

	assert.Equal(t, "", context.Describe(metric))
}

func TestBaseContext_Describe_Format(t *testing.T) {
	context := NewContext("ctx", "name=%<name>s value=%<value>s unit=%<unit>s value_unit=%<value_unit>s")
	metric := NewMetric("metric", 42, "s", nil, "")

	assert.Equal(t, "name=metric value=42 unit=s value_unit=42s", context.Describe(metric))
}

func TestBaseContext_Evaluate(t *testing.T) {
	context := NewContext("ctx", "")
	context.SetResultFactory(NewResultFactory())

	metric := NewMetric("metric", 42, "", nil, "")
	resource := NewResource()
	expected := NewResult(StateOk, metric, context, resource, "")

	assert.Equal(t, expected, context.Evaluate(metric, resource))
}

func TestBaseContext_Performance(t *testing.T) {
	context := NewContext("ctx", "")
	metric := NewMetric("metric", 42, "", nil, "")
	resource := NewResource()

	assert.Nil(t, context.Performance(metric, resource))
}

func TestScalarContext_Evaluate(t *testing.T) {
	warningRange := ParseRange("0:2")
	criticalRange := ParseRange("0:4")
	context := NewScalarContext("ctx", warningRange, criticalRange)

	tests := map[string]struct {
		metricValue float64
		resultState State
		resultHint  string
	}{
		"OK":       {1, StateOk, ""},
		"Warning":  {3, StateWarning, warningRange.ViolationHint()},
		"Critical": {5, StateCritical, criticalRange.ViolationHint()},
	}

	for _, test := range tests {
		metric := NewMetric("metric", test.metricValue, "", nil, "")
		expected := NewResult(test.resultState, metric, context, nil, test.resultHint)

		assert.Equal(t, expected, context.Evaluate(metric, nil))
	}
}

func TestScalarContext_Performance(t *testing.T) {
	warningRange := ParseRange("0:2")
	criticalRange := ParseRange("0:4")
	context := NewScalarContext("ctx", warningRange, criticalRange)

	metric := NewMetric("metric", 42, "", nil, "")
	resource := NewResource()

	expected := &PerfData{
		metric:        metric,
		warningRange:  warningRange,
		criticalRange: criticalRange,
	}

	assert.Equal(t, expected, context.Performance(metric, resource))
}
