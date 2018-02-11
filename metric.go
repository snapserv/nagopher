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
	"strconv"
)

// Metric represents a interface for all metric types.
type Metric interface {
	fmt.Stringer

	Name() string
	Value() float64
	Unit() string
	ValueUnit() string
	ValueRange() *Range
	ContextName() string
}

// BaseMetric represents a generic context from which all other metric types should originate.
type BaseMetric struct {
	name        string
	value       float64
	valueUnit   string
	valueRange  *Range
	contextName string
}

// NewMetric instantiates 'BaseMetric' with the given name, value, unit, range and context name. If no context name is
// given, it is automatically set to the name of the metric.
func NewMetric(name string, value float64, valueUnit string, valueRange *Range, contextName string) *BaseMetric {
	if contextName == "" {
		contextName = name
	}

	return &BaseMetric{
		name:        name,
		value:       value,
		valueUnit:   valueUnit,
		valueRange:  valueRange,
		contextName: contextName,
	}
}

// String is an alias for 'ValueUnit()' to implement the 'fmt.Stringer' interface.
func (m *BaseMetric) String() string {
	return m.ValueUnit()
}

// IsIntegral returns if the current value can be represented as an integer without loosing any precision.
func (m *BaseMetric) IsIntegral() bool {
	return m.value == float64(int64(m.value))
}

// Name represents a getter for the 'name' attribute.
func (m *BaseMetric) Name() string {
	return m.name
}

// Value represents a getter for the 'value' attribute.
func (m *BaseMetric) Value() float64 {
	return m.value
}

// Unit represents a getter for the 'valueUnit' attribute.
func (m *BaseMetric) Unit() string {
	return m.valueUnit
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
	if m.IsIntegral() {
		return fmt.Sprintf("%d%s", int64(m.value), m.valueUnit)
	}

	return fmt.Sprintf("%s%s", strconv.FormatFloat(m.value, 'f', -1, strconv.IntSize), m.valueUnit)
}
