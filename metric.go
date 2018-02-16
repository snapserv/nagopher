/*
 * nagopher - Library for writing Nagios plugins in Go
 * Copyright (C) 2018  Pascal Mathis
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

// Metric represents a interface for all metric types.
type Metric interface {
	fmt.Stringer

	Name() string
	Unit() string
	ValueString() string
	ValueUnit() string
	ValueRange() *Range
	ContextName() string
}

// BaseMetric represents a generic context from which all other metric types should originate.
type BaseMetric struct {
	name        string
	valueUnit   string
	valueRange  *Range
	contextName string
}

// NumericMetric represents a metric type storing numbers as float64.
type NumericMetric struct {
	*BaseMetric
	value float64
}

// StringMetric represents a metric type storing strings.
type StringMetric struct {
	*BaseMetric
	value string
}

// NewBaseMetric instantiates 'BaseMetric' with the given name, unit, range and context name. If no context name is
// given, it is automatically set to the name of the metric.
func NewBaseMetric(name string, valueUnit string, valueRange *Range, contextName string) *BaseMetric {
	if contextName == "" {
		contextName = name
	}

	return &BaseMetric{
		name:        name,
		valueUnit:   valueUnit,
		valueRange:  valueRange,
		contextName: contextName,
	}
}

// String is an alias for 'ValueUnit()' to implement the 'fmt.Stringer' interface.
func (m *BaseMetric) String() string {
	return m.ValueUnit()
}

// Name represents a getter for the 'name' attribute.
func (m *BaseMetric) Name() string {
	return m.name
}

// Unit represents a getter for the 'valueUnit' attribute.
func (m *BaseMetric) Unit() string {
	return m.valueUnit
}

// ValueString returns the string representation of the metric value
func (m *BaseMetric) ValueString() string {
	return ""
}

// ValueRange represents a getter for the 'valueRange' attribute.
func (m *BaseMetric) ValueRange() *Range {
	return m.valueRange
}

// ContextName represents a getter for the 'contextName' attribute.
func (m *BaseMetric) ContextName() string {
	return m.contextName
}

// ValueUnit returns the string representation of the metric, which is '<value><unit>'.
func (m *BaseMetric) ValueUnit() string {
	return "N/A"
}

// NewNumericMetric instantiates 'NumericMetric', which represents a metric containing a numeric value. All numbers will
// be stored as float64 internally, however always the shortest value without losing precision will be displayed.
func NewNumericMetric(name string, value float64, valueUnit string, valueRange *Range, contextName string) *NumericMetric {
	return &NumericMetric{
		BaseMetric: NewBaseMetric(name, valueUnit, valueRange, contextName),
		value:      value,
	}
}

// IsIntegral returns if the current value can be represented as an integer without loosing any precision.
func (m *NumericMetric) IsIntegral() bool {
	return m.value == float64(int64(m.value))
}

// Value represents a getter for the 'value' attribute.
func (m *NumericMetric) Value() float64 {
	return m.value
}

// ValueString returns the string representation of the metric value
func (m *NumericMetric) ValueString() string {
	if m.IsIntegral() {
		return fmt.Sprintf("%d", int64(m.value))
	}

	return strconv.FormatFloat(m.value, 'f', -1, strconv.IntSize)
}

// ValueUnit returns the string representation of the metric, which is '<value><unit>'.
func (m *NumericMetric) ValueUnit() string {
	if math.IsNaN(m.value) {
		return "U"
	}

	return m.ValueString() + m.Unit()
}

// NewStringMetric instantiates 'StringMetric', which represents a metric containing a string value. Please note that
// this metric type does not support any unit or value ranges, for which NumericMetric should be used instead.
func NewStringMetric(name string, value string, contextName string) *StringMetric {
	return &StringMetric{
		BaseMetric: NewBaseMetric(name, "", nil, contextName),
		value:      value,
	}
}

// Value represents a getter for the 'value' attribute.
func (m *StringMetric) Value() string {
	return m.value
}

// ValueString returns the string representation of the metric value
func (m *StringMetric) ValueString() string {
	return m.value
}

// ValueUnit returns the string representation of the metric, which is '<value><unit>'.
func (m *StringMetric) ValueUnit() string {
	return m.ValueString() + m.Unit()
}
