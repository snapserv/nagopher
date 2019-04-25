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
	"strings"
)

type stringInfoContext struct {
	baseContext
}

type stringMatchContext struct {
	baseContext

	problemState   State
	expectedValues []string
}

// NewStringInfoContext instantiates a Context which only holds and returns a plain string without any further logic.
func NewStringInfoContext(name string) Context {
	stringInfoContext := &stringInfoContext{
		baseContext: *newBaseContext(name, "%<value>s"),
	}

	return stringInfoContext
}

func (c stringInfoContext) Evaluate(metric Metric, resource Resource) Result {
	return NewResult(
		ResultState(StateOk()),
		ResultMetric(metric), ResultContext(c), ResultResource(resource),
	)
}

// NewStringMatchContext instantiates a Context which holds a string, which is being compared to a whitelist of
// acceptable values during the evaluation phase. Should the value not be accepted, a problem state gets returned.
func NewStringMatchContext(name string, problemState State, expectedValues []string) Context {
	stringContext := &stringMatchContext{
		baseContext: *newBaseContext(name, "%<name>s is %<value>s"),

		problemState:   problemState,
		expectedValues: stringsToLower(expectedValues),
	}

	return stringContext
}

func (c stringMatchContext) Evaluate(metric Metric, resource Resource) Result {
	stringMetric, ok := metric.(StringMetric)
	if !ok {
		return NewResult(
			ResultState(StateUnknown()),
			ResultMetric(metric), ResultContext(c), ResultResource(resource),
			ResultHint(fmt.Sprintf("StringMatchContext can not process metric of type [%s]", reflect.TypeOf(metric))),
		)
	}

	if len(c.expectedValues) == 0 {
		return NewResult(
			ResultState(StateOk()),
			ResultMetric(metric), ResultContext(c), ResultResource(resource),
		)
	}

	value := strings.ToLower(stringMetric.Value())
	for _, expectedValue := range c.expectedValues {
		if value == expectedValue {
			return NewResult(
				ResultState(StateOk()),
				ResultMetric(metric), ResultContext(c), ResultResource(resource),
			)
		}
	}

	return NewResult(
		ResultState(c.problemState),
		ResultMetric(metric), ResultContext(c), ResultResource(resource),
		ResultHint(fmt.Sprintf("got [%s], expected [%s]", value, strings.Join(c.expectedValues, "],["))),
	)
}

func stringsToLower(values []string) []string {
	var sanitizedValues []string

	for _, value := range values {
		sanitizedValues = append(sanitizedValues, strings.ToLower(value))
	}

	return sanitizedValues
}
