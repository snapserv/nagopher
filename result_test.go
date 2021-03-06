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

func TestNewResult(t *testing.T) {
	// given
	expectedState := StateOk()
	expectedMetric := MustNewStringMetric("metric", "test", "context")
	expectedContext := NewStringInfoContext("context")
	expectedResource := NewResource()
	expectedHint := "Result Hint"

	// when
	result := NewResult(
		ResultState(expectedState), ResultHint(expectedHint),
		ResultMetric(expectedMetric), ResultContext(expectedContext), ResultResource(expectedResource),
	)
	actualState, _ := result.State().Get()
	actualMetric, _ := result.Metric().Get()
	actualContext, _ := result.Context().Get()
	actualResource, _ := result.Resource().Get()

	// then
	assert.Equal(t, expectedHint, result.Hint())
	assert.Equal(t, expectedState, actualState)
	assert.Equal(t, expectedMetric, actualMetric)
	assert.Equal(t, expectedContext, actualContext)
	assert.Equal(t, expectedResource, actualResource)
}

func TestResult_String(t *testing.T) {
	// given
	metric := MustNewNumericMetric("metric", 13.37, "", nil, "")
	context := NewScalarContext("context", nil, nil)
	hint := "Result Hint"

	// when
	result1 := NewResult(ResultContext(context), ResultMetric(metric))
	result2 := NewResult(ResultMetric(metric))
	result3 := NewResult(ResultContext(context), ResultMetric(metric), ResultHint(hint))
	result4 := NewResult(ResultMetric(metric), ResultHint(hint))

	// then
	assert.Equal(t, "metric is 13.37", result1.String())
	assert.Equal(t, "13.37", result2.String())
	assert.Equal(t, "metric is 13.37 (Result Hint)", result3.String())
	assert.Equal(t, "13.37 (Result Hint)", result4.String())
}
