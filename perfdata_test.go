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

func TestNewPerfData(t *testing.T) {
	// given
	expectedMetric1 := MustNewStringMetric("test", "Hello World", "")
	expectedMetric2 := MustNewStringMetric("inv'=alid", "Hello World", "")

	// when
	perfData1, err1 := NewPerfData(expectedMetric1, nil, nil)
	perfData2, err2 := NewPerfData(expectedMetric2, nil, nil)

	// then
	assert.NoError(t, err1)
	assert.Error(t, err2)
	assert.Implements(t, (*PerfData)(nil), perfData1)
	assert.Equal(t, expectedMetric1, perfData1.Metric())
	assert.Nil(t, perfData2)
}

func TestNewNumericPerfData(t *testing.T) {
	// when
	perfData1, err1 := NewNumericPerfData("test", 13.37, "", nil, nil, nil)
	perfData2, err2 := NewNumericPerfData("inv'=alid", 13.37, "", nil, nil, nil)
	perfData3, err3 := NewNumericPerfData("", 13.37, "", nil, nil, nil)

	// then
	assert.NoError(t, err1)
	assert.Error(t, err2)
	assert.Error(t, err3)

	assert.Implements(t, (*PerfData)(nil), perfData1)
	assert.Nil(t, perfData2)
	assert.Nil(t, perfData3)
}

func TestPerfData_ToNagiosPerfData(t *testing.T) {
	// given
	valueRange := NewBounds(LowerBound(-100), UpperBound(100))
	warningThreshold := NewBounds(LowerBound(10), UpperBound(20), InvertedBounds(true))
	criticalThreshold := NewBounds(LowerBound(10), UpperBound(20))
	metric1 := MustNewNumericMetric("test", 13.37, "B", &valueRange, "")
	metric2 := MustNewNumericMetric("test with quoting", 42, "X", nil, "")

	// when
	perfData1, err := NewPerfData(metric1, &warningThreshold, &criticalThreshold)
	perfData2, err := NewPerfData(metric2, nil, nil)

	// then
	assert.NoError(t, err)
	assert.Equal(t, "test=13.37B;@10:20;10:20;-100;100", perfData1.ToNagiosPerfData())
	assert.Equal(t, "'test with quoting'=42X", perfData2.ToNagiosPerfData())
}
