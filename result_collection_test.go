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

func TestResultCollection_Add(t *testing.T) {
	// given
	result1 := NewResult()
	result2 := NewResult()
	result3 := NewResult()
	results := NewResultCollection()

	// when
	results.Add(result1)
	results.Add(result2, result3)

	// then
	assert.Equal(t, 3, results.Count())
	assert.Contains(t, results.Get(), result1)
	assert.Contains(t, results.Get(), result2)
	assert.Contains(t, results.Get(), result3)
}

func TestResultCollection_MostSignificantState(t *testing.T) {
	// given
	results1 := NewResultCollection()
	results2 := NewResultCollection()

	results1.Add(
		NewResult(ResultState(StateOk())),
		NewResult(),
		NewResult(ResultState(StateCritical())),
		NewResult(),
	)

	// when
	state1, err1 := results1.MostSignificantState().Get()
	state2, err2 := results2.MostSignificantState().Get()

	// then
	assert.NoError(t, err1)
	assert.Error(t, err2)
	assert.Equal(t, StateCritical(), state1)
	assert.Nil(t, state2)
}

func TestResultCollection_MostSignificantState2(t *testing.T) {
	// given
	results := []ResultCollection{
		NewResultCollection(), NewResultCollection(), NewResultCollection(),
		NewResultCollection(), NewResultCollection(), NewResultCollection(),
	}

	resultInfo := NewResult(ResultState(StateInfo()))
	resultOk := NewResult(ResultState(StateOk()))
	resultWarning := NewResult(ResultState(StateWarning()))
	resultCritical := NewResult(ResultState(StateCritical()))

	// when
	results[1].Add(resultInfo)
	results[2].Add(resultInfo, resultOk)
	results[3].Add(resultOk, resultInfo)
	results[4].Add(resultInfo, resultOk, resultWarning)
	results[5].Add(resultInfo, resultOk, resultWarning, resultCritical)

	// then
	assert.Nil(t, results[0].MostSignificantState().OrElse(nil))
	assert.Equal(t, StateInfo(), results[1].MostSignificantState().OrElse(nil))
	assert.Equal(t, StateOk(), results[2].MostSignificantState().OrElse(nil))
	assert.Equal(t, StateOk(), results[3].MostSignificantState().OrElse(nil))
	assert.Equal(t, StateWarning(), results[4].MostSignificantState().OrElse(nil))
	assert.Equal(t, StateCritical(), results[5].MostSignificantState().OrElse(nil))
}

func TestResultCollection_GetMetricByName(t *testing.T) {
	// given
	expectedMetric1 := MustNewStringMetric("metric 1", "Hello", "")
	expectedMetric2 := MustNewNumericMetric("metric 2", 13.37, "", nil, "")
	results := NewResultCollection()
	results.Add(
		NewResult(),
		NewResult(ResultMetric(expectedMetric1)),
		NewResult(ResultMetric(expectedMetric2)),
	)

	// when
	actualMetric1, err1 := results.GetMetricByName("metric 1").Get()
	actualMetric2, err2 := results.GetMetricByName("metric 2").Get()
	actualMetric3, err3 := results.GetMetricByName("missing").Get()

	// then
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Error(t, err3)
	assert.Equal(t, expectedMetric1, actualMetric1)
	assert.Equal(t, expectedMetric2, actualMetric2)
	assert.Nil(t, actualMetric3)
}

func TestResultCollection_GetNumericMetricValue(t *testing.T) {
	// given
	expectedMetric := MustNewNumericMetric("metric", 13.37, "", nil, "")
	results := NewResultCollection()
	results.Add(NewResult(ResultMetric(expectedMetric)))

	// when
	actualMetric1, err1 := results.GetNumericMetricValue("metric").Get()
	actualMetric2, err2 := results.GetNumericMetricValue("missing").Get()

	// then
	assert.NoError(t, err1)
	assert.Error(t, err2)
	assert.Equal(t, expectedMetric.Value(), actualMetric1)
	assert.Equal(t, float64(0), actualMetric2)
}

func TestResultCollection_GetStringMetricValue(t *testing.T) {
	// given
	expectedMetric := MustNewStringMetric("metric", "Hello World", "")
	results := NewResultCollection()
	results.Add(NewResult(ResultMetric(expectedMetric)))

	// when
	actualMetric1, err1 := results.GetStringMetricValue("metric").Get()
	actualMetric2, err2 := results.GetStringMetricValue("missing").Get()

	// then
	assert.NoError(t, err1)
	assert.Error(t, err2)
	assert.Equal(t, expectedMetric.Value(), actualMetric1)
	assert.Equal(t, "", actualMetric2)
}
