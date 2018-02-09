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

import "github.com/chonla/format"

type Context interface {
	Name() string
	Describe(Metric) string
	Evaluate(Metric, Resource) Result
	Performance(Metric, Resource) *PerfData
}

type BaseContext struct {
	name          string
	format        string
	resultFactory ResultFactory
}

type ScalarContext struct {
	*BaseContext

	warningRange  *Range
	criticalRange *Range
}

func NewContext(name string, format string) *BaseContext {
	return &BaseContext{
		name:          name,
		format:        format,
		resultFactory: NewResultFactory(),
	}
}

func (c *BaseContext) SetResultFactory(resultFactory ResultFactory) {
	c.resultFactory = resultFactory
}

func (c *BaseContext) Name() string {
	return c.name
}

func (c *BaseContext) Describe(metric Metric) string {
	data := map[string]interface{}{
		"name":       metric.Name(),
		"value":      metric.Value(),
		"unit":       metric.Unit(),
		"value_unit": metric.ValueUnit(),
	}

	return format.Sprintf(c.format, data)
}

func (c *BaseContext) Evaluate(metric Metric, resource Resource) Result {
	return c.resultFactory(StateOk, metric, c, resource, "")
}

func (c *BaseContext) Performance(metric Metric, resource Resource) *PerfData {
	return nil
}

func NewScalarContext(name string, warningRange *Range, criticalRange *Range) *ScalarContext {
	return &ScalarContext{
		BaseContext: NewContext(name, "%<name>s is %<value_unit>s"),

		warningRange:  warningRange,
		criticalRange: criticalRange,
	}
}

func (c *ScalarContext) Evaluate(metric Metric, resource Resource) Result {
	if c.criticalRange != nil && !c.criticalRange.Match(metric.Value()) {
		return c.resultFactory(StateCritical, metric, c, resource, c.criticalRange.ViolationHint())
	} else if c.warningRange != nil && !c.warningRange.Match(metric.Value()) {
		return c.resultFactory(StateWarning, metric, c, resource, c.warningRange.ViolationHint())
	} else {
		return c.resultFactory(StateOk, metric, c, resource, "")
	}
}

func (c *ScalarContext) Performance(metric Metric, resource Resource) *PerfData {
	return &PerfData{
		metric:        metric,
		warningRange:  c.warningRange,
		criticalRange: c.criticalRange,
	}
}
