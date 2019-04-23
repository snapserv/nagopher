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
	"math"
	"strconv"
)

// NumericMetric represents a Metric storing float64 values
type NumericMetric interface {
	Metric

	Value() float64
}

type numericMetric struct {
	baseMetric
	value float64
}

// NewNumericMetric instantiates a new NumericMetric with the given parameters.
func NewNumericMetric(name string, value float64, valueUnit string, valueRange *Bounds, contextName string) (NumericMetric, error) {
	baseMetric, err := newBaseMetric(name, valueUnit, valueRange, contextName)
	if err != nil {
		return nil, err
	}

	numericMetric := &numericMetric{
		baseMetric: *baseMetric,
		value:      value,
	}

	return numericMetric, nil
}

// MustNewNumericMetric calls MustNewNumericMetric and panics in case the creation of a metric instance fails
func MustNewNumericMetric(name string, value float64, valueUnit string, valueRange *Bounds, contextName string) NumericMetric {
	metric, err := NewNumericMetric(name, value, valueUnit, valueRange, contextName)
	if err != nil {
		panic(err)
	}

	return metric
}

func (m numericMetric) ToNagiosValue() string {
	if math.IsNaN(m.value) {
		return "U"
	}

	return m.ValueString() + m.ValueUnit()
}

func (m numericMetric) Value() float64 {
	return m.value
}

func (m numericMetric) ValueString() string {
	if m.IsIntegral() {
		return fmt.Sprintf("%d", int64(m.value))
	}

	return strconv.FormatFloat(m.value, 'f', -1, strconv.IntSize)
}

func (m numericMetric) IsIntegral() bool {
	return m.value == float64(int64(m.value))
}
