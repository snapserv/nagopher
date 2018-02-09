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

type Metric interface {
	fmt.Stringer

	Name() string
	Value() float64
	Unit() string
	ValueUnit() string
	ValueRange() *Range
	ContextName() string
}

type BaseMetric struct {
	name        string
	value       float64
	valueUnit   string
	valueRange  *Range
	contextName string
}

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

func (m *BaseMetric) String() string {
	return m.ValueUnit()
}

func (m *BaseMetric) IsIntegral() bool {
	return m.value == float64(int64(m.value))
}

func (m *BaseMetric) Name() string {
	return m.name
}

func (m *BaseMetric) Value() float64 {
	return m.value
}

func (m *BaseMetric) Unit() string {
	return m.valueUnit
}

func (m *BaseMetric) ValueUnit() string {
	if m.IsIntegral() {
		return fmt.Sprintf("%d%s", int64(m.value), m.valueUnit)
	} else {
		return fmt.Sprintf("%s%s",
			strconv.FormatFloat(m.value, 'f', -1, strconv.IntSize),
			m.valueUnit)
	}
}

func (m *BaseMetric) ValueRange() *Range {
	return m.valueRange
}

func (m *BaseMetric) ContextName() string {
	return m.contextName
}
