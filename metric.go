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

//go:generate optional -type=Metric

import (
	"errors"
	"github.com/markphelps/optional"
)

// Metric stores a value of a given type associated to a specific context, optionally restrained into a specific range.
// It can be converted to the according representation of that value in context of Nagios plugins.
type Metric interface {
	ToNagiosValue() string

	Name() string
	ValueUnit() string
	ValueString() string
	ValueRange() OptionalBounds
	ContextName() string
}

type baseMetric struct {
	name        string
	valueUnit   string
	valueRange  OptionalBounds
	contextName optional.String
}

func newBaseMetric(name string, valueUnit string, valueRange *Bounds, contextName string) (*baseMetric, error) {
	if name == "" {
		return nil, errors.New("metric name must not be empty")
	}

	baseMetric := &baseMetric{
		name:      name,
		valueUnit: valueUnit,
	}

	if valueRange != nil {
		if (*valueRange).IsInverted() {
			return nil, errors.New("metric value range must not be inverted")
		}

		baseMetric.valueRange = NewOptionalBounds(*valueRange)
	}
	if contextName != "" {
		baseMetric.contextName = optional.NewString(contextName)
	}

	return baseMetric, nil
}

func (m baseMetric) Name() string {
	return m.name
}

func (m baseMetric) ValueUnit() string {
	return m.valueUnit
}

func (m baseMetric) ValueRange() OptionalBounds {
	return m.valueRange
}

func (m baseMetric) ContextName() string {
	return m.contextName.OrElse(m.name)
}
