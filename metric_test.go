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

func TestNewBaseMetric(t *testing.T) {
	// given
	expectedValueRange := NewBounds()

	// when
	metric, err := newBaseMetric("metric", "B", &expectedValueRange, "context")

	// then
	assert.NoError(t, err)
	assert.Equal(t, "metric", metric.Name())
	assert.Equal(t, "B", metric.ValueUnit())
	assert.Equal(t, "context", metric.ContextName())

	actualValueRange, err := metric.ValueRange().Get()
	assert.NoError(t, err)
	assert.Equal(t, expectedValueRange, actualValueRange)
}

func TestNewBaseMetric_Invalid(t *testing.T) {
	// given
	invertedBounds := NewBounds(InvertedBounds(true))

	// when
	metric1, err1 := newBaseMetric("", "", nil, "")
	metric2, err2 := newBaseMetric("metric", "", &invertedBounds, "")

	// then
	assert.Error(t, err1)
	assert.Error(t, err2)
	assert.Nil(t, metric1)
	assert.Nil(t, metric2)
}

func TestBaseMetric_ContextName_Empty(t *testing.T) {
	// when
	metric, err := newBaseMetric("metric", "", nil, "")

	// then
	assert.NoError(t, err)
	assert.Equal(t, "metric", metric.ContextName())
}
