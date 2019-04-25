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

func TestBaseContext_Describe(t *testing.T) {
	// given
	metric := MustNewNumericMetric("test", 13.37, "apples", nil, "")
	context := newBaseContext("Test Context", "name=%<name>s value=%<value>s unit=%<unit>s")

	// when
	description := context.Describe(metric)

	// then
	assert.Equal(t, "name=test value=13.37 unit=apples", description)
}

func TestBaseContext_Evaluate(t *testing.T) {
	// given
	context := newBaseContext("Test Context", "%<value>s")
	expectedMetric := MustNewStringMetric("string", "Hello", "")
	expectedResource := NewResource()

	// when
	result := context.Evaluate(expectedMetric, expectedResource)
	resultState, _ := result.State().Get()
	resultContext, _ := result.Context().Get()
	resultMetric, _ := result.Metric().Get()
	resultResource, _ := result.Resource().Get()

	// then
	assert.Equal(t, true, result.State().Present())
	assert.Equal(t, true, result.Context().Present())
	assert.Equal(t, true, result.Metric().Present())
	assert.Equal(t, true, result.Resource().Present())

	assert.Equal(t, StateOk(), resultState)
	assert.Equal(t, *context, resultContext)
	assert.Equal(t, expectedMetric, resultMetric)
	assert.Equal(t, expectedResource, resultResource)
}

func TestBaseContext_Performance(t *testing.T) {
	// given
	context := newBaseContext("Test Context", "%<value>s")
	metric := MustNewStringMetric("string", "Hello", "")
	resource := NewResource()

	// when
	perfData, err := context.Performance(metric, resource)

	// then
	assert.NoError(t, err)
	assert.Equal(t, OptionalPerfData{}, perfData)
}
