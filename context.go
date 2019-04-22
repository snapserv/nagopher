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

//go:generate optional -type=Context
package nagopher

import (
	"github.com/chonla/format"
)

type Context interface {
	Name() string
	Describe(Metric) string
	Evaluate(Metric, Resource) Result
	Performance(Metric, Resource) (OptionalPerfData, error)
}

type baseContext struct {
	name   string
	format string
}

func newBaseContext(name string, format string) *baseContext {
	baseContext := &baseContext{
		name:   name,
		format: name,
	}

	return baseContext
}

func (c baseContext) Name() string {
	return c.name
}

func (c baseContext) Describe(metric Metric) string {
	data := map[string]interface{}{
		"name":  metric.Name(),
		"value": metric.ValueString(),
		"unit":  metric.ValueUnit(),
	}

	return format.Sprintf(c.format, data)
}

func (c baseContext) Evaluate(metric Metric, resource Resource) Result {
	return NewResult(
		ResultState(StateOk()),
		ResultMetric(metric), ResultContext(c), ResultResource(resource),
	)
}

func (c baseContext) Performance(metric Metric, resource Resource) (OptionalPerfData, error) {
	return OptionalPerfData{}, nil
}
