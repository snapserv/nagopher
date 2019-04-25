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

func TestOptionalMetric_Empty(t *testing.T) {
	// when
	optionalMetric := OptionalMetric{}
	actualMetric, err := optionalMetric.Get()

	// then
	assert.Error(t, err)
	assert.Empty(t, actualMetric)
}

func TestNewOptionalMetric(t *testing.T) {
	// given
	expectedMetric, _ := NewStringMetric("expected", "", "")

	// when
	optionalMetric := NewOptionalMetric(expectedMetric)
	actualMetric, err := optionalMetric.Get()

	// then
	assert.NoError(t, err)
	assert.Equal(t, true, optionalMetric.Present())
	assert.Equal(t, expectedMetric, actualMetric)
}

func TestOptionalMetric_OrElse(t *testing.T) {
	// given
	expectedMetric, _ := NewStringMetric("expected", "", "")
	alternativeMetric, _ := NewStringMetric("alternative", "", "")

	// when
	optionalMetric1 := NewOptionalMetric(expectedMetric)
	optionalMetric2 := OptionalMetric{}

	// then
	assert.Equal(t, expectedMetric, optionalMetric1.OrElse(alternativeMetric))
	assert.Equal(t, alternativeMetric, optionalMetric2.OrElse(alternativeMetric))
}

func TestOptionalMetric_Set(t *testing.T) {
	// given
	metric1, _ := NewStringMetric("metric 1", "", "")
	metric2, _ := NewStringMetric("metric 2", "", "")
	alternativeMetric, _ := NewStringMetric("alternative", "", "")

	// when
	optionalMetric := NewOptionalMetric(metric1)
	optionalMetric.Set(metric2)

	// then
	assert.Equal(t, true, optionalMetric.Present())
	assert.Equal(t, metric2, optionalMetric.OrElse(alternativeMetric))
}

func TestOptionalMetric_If(t *testing.T) {
	// given
	var actualMetricPtr *Metric
	var expectedMetric Metric
	expectedMetric, _ = NewStringMetric("expected", "", "")

	// when
	optionalMetric := NewOptionalMetric(expectedMetric)
	optionalMetric.If(func(metric Metric) {
		actualMetricPtr = &metric
	})

	// then
	assert.Equal(t, &expectedMetric, actualMetricPtr)
}
