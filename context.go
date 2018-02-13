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
	"reflect"
	"strconv"

	"github.com/chonla/format"
)

// Context represents a interface for all context types.
type Context interface {
	Name() string
	Describe(Metric) string
	Evaluate(Metric, Resource) Result
	Performance(Metric, Resource) *PerfData
}

// BaseContext represents a generic context from which all other context types should originate.
type BaseContext struct {
	name          string
	format        string
	resultFactory ResultFactory
}

// ScalarContext represents a context for scalar values with support for warning/critical threshold ranges.
type ScalarContext struct {
	*BaseContext

	warningRange  *Range
	criticalRange *Range
}

// NewContext instantiates 'BaseContext' with a given name and format string.
func NewContext(name string, format string) *BaseContext {
	return &BaseContext{
		name:          name,
		format:        format,
		resultFactory: NewResultFactory(),
	}
}

// SetResultFactory allows overriding the default result factory, which by default instantiates 'BaseResult'.
func (c *BaseContext) SetResultFactory(resultFactory ResultFactory) {
	c.resultFactory = resultFactory
}

// Name represents a getter for the 'name' attribute.
func (c *BaseContext) Name() string {
	return c.name
}

// Describe returns a formatted string based on the 'format' attribute for the given metric.
func (c *BaseContext) Describe(metric Metric) string {
	var data map[string]interface{}

	switch m := metric.(type) {
	case *NumericMetric:
		data = map[string]interface{}{
			"name":       metric.Name(),
			"value":      strconv.FormatFloat(m.Value(), 'f', -1, strconv.IntSize),
			"unit":       metric.Unit(),
			"value_unit": metric.ValueUnit(),
		}

	case *StringMetric:
		data = map[string]interface{}{
			"name":       metric.Name(),
			"value":      m.Value(),
			"unit":       metric.Unit(),
			"value_unit": metric.ValueUnit(),
		}
	}

	return format.Sprintf(c.format, data)
}

// Evaluate returns a Result object for the given metric and resource.
func (c *BaseContext) Evaluate(metric Metric, resource Resource) Result {
	return c.resultFactory(StateOk, metric, c, resource, "")
}

// Performance returns performance data for the given metric and resource.
func (c *BaseContext) Performance(metric Metric, resource Resource) *PerfData {
	return nil
}

// NewScalarContext instantiates 'ScalarContext' with the given name and optional warning/critical threshold ranges.
func NewScalarContext(name string, warningRange *Range, criticalRange *Range) *ScalarContext {
	return &ScalarContext{
		BaseContext: NewContext(name, "%<name>s is %<value_unit>s"),

		warningRange:  warningRange,
		criticalRange: criticalRange,
	}
}

// Evaluate checks if the given metric and resource match the warning/critical threshold ranges, if given,
// and returns the appropriate Result object.
func (c *ScalarContext) Evaluate(metric Metric, resource Resource) Result {
	numberMetric, ok := metric.(*NumericMetric)
	if !ok {
		return c.resultFactory(StateUnknown, metric, c, resource,
			fmt.Sprintf("ScalarContext can not process metrics of type [%s]", reflect.TypeOf(metric)))
	}

	if c.criticalRange != nil && !c.criticalRange.Match(numberMetric.Value()) {
		return c.resultFactory(StateCritical, metric, c, resource, c.criticalRange.ViolationHint())
	} else if c.warningRange != nil && !c.warningRange.Match(numberMetric.Value()) {
		return c.resultFactory(StateWarning, metric, c, resource, c.warningRange.ViolationHint())
	} else {
		return c.resultFactory(StateOk, metric, c, resource, "")
	}
}

// Performance returns the performance data for the scalar context, including the threshold ranges if set.
func (c *ScalarContext) Performance(metric Metric, resource Resource) *PerfData {
	return &PerfData{
		metric:        metric,
		warningRange:  c.warningRange,
		criticalRange: c.criticalRange,
	}
}
