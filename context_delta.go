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
	"reflect"
)

type deltaContext struct {
	scalarContext

	previousValue *float64
}

// NewDeltaContext creates a new scalar Context object, which operates the same way as a ScalarContext, but instead
// of using the current absolute metric value, it will be compared to a previous measurement. It is the callers duty
// to provide a pointer to the previous metric value or nil, if not available.
func NewDeltaContext(name string, previousValue *float64, warningThreshold *Bounds, criticalThreshold *Bounds) Context {
	baseContext := NewScalarContext(name, warningThreshold, criticalThreshold)
	scalarContext := baseContext.(*scalarContext)
	deltaContext := &deltaContext{
		scalarContext: *scalarContext,
		previousValue: previousValue,
	}

	return deltaContext
}

func (c *deltaContext) Evaluate(metric Metric, resource Resource) Result {
	numericMetric, ok := metric.(NumericMetric)
	if !ok {
		return NewResult(
			ResultState(StateUnknown()),
			ResultMetric(metric), ResultContext(c), ResultResource(resource),
			ResultHint(fmt.Sprintf("DeltaContext can not process metric of type [%s]", reflect.TypeOf(metric))),
		)
	}

	metricValue := numericMetric.Value()
	previousValue := float64(0)
	if c.previousValue != nil {
		previousValue = *c.previousValue
		*c.previousValue = metricValue
	}

	deltaValue := metricValue - previousValue
	deltaMetric := MustNewNumericMetric(numericMetric.Name()+"_delta", deltaValue, "", nil, numericMetric.ContextName())

	emptyBounds := NewBounds()
	warningThreshold := c.warningThreshold.OrElse(emptyBounds)
	criticalThreshold := c.criticalThreshold.OrElse(emptyBounds)

	if !criticalThreshold.Match(deltaValue) {
		return NewResult(
			ResultState(StateCritical()),
			ResultMetric(deltaMetric), ResultContext(c), ResultResource(resource),
			ResultHint(criticalThreshold.ViolationHint()),
		)
	} else if !warningThreshold.Match(deltaValue) {
		return NewResult(
			ResultState(StateWarning()),
			ResultMetric(deltaMetric), ResultContext(c), ResultResource(resource),
			ResultHint(warningThreshold.ViolationHint()),
		)
	}

	return NewResult(
		ResultState(StateOk()),
		ResultMetric(deltaMetric), ResultContext(c), ResultResource(resource),
	)
}

func (c deltaContext) Performance(metric Metric, resource Resource) (OptionalPerfData, error) {
	perfData, err := NewPerfData(metric, nil, nil)
	if err != nil {
		return OptionalPerfData{}, err
	}

	return NewOptionalPerfData(perfData), nil
}
