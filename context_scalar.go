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

type scalarContext struct {
	baseContext

	warningThreshold  OptionalBounds
	criticalThreshold OptionalBounds
}

func NewScalarContext(name string, warningThreshold *Bounds, criticalThreshold *Bounds) Context {
	scalarContext := &scalarContext{
		baseContext: *newBaseContext(name, "%<name>s is %<value>s%<unit>s"),
	}

	if warningThreshold != nil {
		scalarContext.warningThreshold = NewOptionalBounds(*warningThreshold)
	}
	if criticalThreshold != nil {
		scalarContext.criticalThreshold = NewOptionalBounds(*criticalThreshold)
	}

	return scalarContext
}

func (c scalarContext) Evaluate(metric Metric, resource Resource) Result {
	numericMetric, ok := metric.(NumericMetric)
	if !ok {
		return NewResult(
			ResultState(StateUnknown()),
			ResultMetric(metric), ResultContext(c), ResultResource(resource),
			ResultHint(fmt.Sprintf("ScalarContext can not process metric of type [%s]", reflect.TypeOf(metric))),
		)
	}

	emptyBounds := NewBounds()
	warningThreshold := c.warningThreshold.OrElse(emptyBounds)
	criticalThreshold := c.criticalThreshold.OrElse(emptyBounds)

	if criticalThreshold.Match(numericMetric.Value()) {
		return NewResult(
			ResultState(StateCritical()),
			ResultMetric(metric), ResultContext(c), ResultResource(resource),
			ResultHint(criticalThreshold.String()),
		)
	} else if warningThreshold.Match(numericMetric.Value()) {
		return NewResult(
			ResultState(StateWarning()),
			ResultMetric(metric), ResultContext(c), ResultResource(resource),
			ResultHint(warningThreshold.String()),
		)
	}

	return NewResult(
		ResultState(StateOk()),
		ResultMetric(metric), ResultContext(c), ResultResource(resource),
	)
}

func (c scalarContext) Performance(metric Metric, resource Resource) (OptionalPerfData, error) {
	var warningThreshold *Bounds
	var criticalThreshold *Bounds

	if threshold, err := c.warningThreshold.Get(); err == nil {
		warningThreshold = &threshold
	}
	if threshold, err := c.criticalThreshold.Get(); err == nil {
		criticalThreshold = &threshold
	}

	perfData, err := NewPerfData(metric, warningThreshold, criticalThreshold)
	if err != nil {
		return OptionalPerfData{}, err
	}

	return NewOptionalPerfData(perfData), nil
}
