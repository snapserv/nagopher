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

func TestDeltaContext_Evaluate(t *testing.T) {
	// given
	previousValue := float64(100)
	warningThreshold := NewBounds(LowerBound(-5), UpperBound(5))
	criticalThreshold := NewBounds(LowerBound(-10), UpperBound(10))
	context := NewDeltaContext("context", &previousValue, &warningThreshold, &criticalThreshold)
	metric1 := MustNewNumericMetric("metric 1", previousValue+3, "", nil, "")
	metric2 := MustNewNumericMetric("metric 2", previousValue+3-7, "", nil, "")
	metric3 := MustNewNumericMetric("metric 3", previousValue+3-7+11, "", nil, "")
	metric4 := MustNewStringMetric("invalid", "Oops!", "")
	resource := NewResource()

	// when
	result1 := context.Evaluate(metric1, resource)
	result2 := context.Evaluate(metric2, resource)
	result3 := context.Evaluate(metric3, resource)
	result4 := context.Evaluate(metric4, resource)

	// then
	assert.Equal(t, StateOk(), result1.State().OrElse(nil))
	assert.Equal(t, StateWarning(), result2.State().OrElse(nil))
	assert.Equal(t, StateCritical(), result3.State().OrElse(nil))
	assert.Equal(t, StateUnknown(), result4.State().OrElse(nil))

	assert.Equal(t, "", result1.Hint())
	assert.Equal(t, "outside range -5:5", result2.Hint())
	assert.Equal(t, "outside range -10:10", result3.Hint())
	assert.Contains(t, result4.Hint(), "DeltaContext can not process metric of type")
}

func TestDeltaContext_Evaluate_Adjust(t *testing.T) {
	// given
	previousValue := float64(100)
	warningThreshold := NewBounds(LowerBound(-5), UpperBound(5))
	context := NewDeltaContext("context", &previousValue, &warningThreshold, nil)
	metric1 := MustNewNumericMetric("metric 1", previousValue+4, "", nil, "")
	metric2 := MustNewNumericMetric("metric 2", previousValue+4+4, "", nil, "")
	resource := NewResource()

	// when
	result1 := context.Evaluate(metric1, resource)
	result2 := context.Evaluate(metric2, resource)

	// then
	assert.Equal(t, StateOk(), result1.State().OrElse(nil))
	assert.Equal(t, StateOk(), result2.State().OrElse(nil))
	assert.Equal(t, "", result1.Hint())
	assert.Equal(t, "", result2.Hint())
}

func TestDeltaContext_Performance(t *testing.T) {
	// given
	previousValue := float64(100)
	context := NewDeltaContext("context", &previousValue, nil, nil)
	metric1 := MustNewNumericMetric("valid", 42, "", nil, "")
	metric2 := MustNewNumericMetric("inv='alid", 42, "", nil, "")
	resource := NewResource()

	// when
	var perfData1, perfData2 PerfData = nil, nil
	optionalPerfData1, err1 := context.Performance(metric1, resource)
	optionalPerfData2, err2 := context.Performance(metric2, resource)

	if err1 == nil {
		perfData1, err1 = optionalPerfData1.Get()
	}
	if err2 == nil {
		perfData2, err2 = optionalPerfData2.Get()
	}

	// then
	assert.NoError(t, err1)
	assert.Error(t, err2)

	assert.Implements(t, (*PerfData)(nil), perfData1)
	assert.Nil(t, perfData2)
}
