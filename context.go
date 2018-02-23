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
	"github.com/chonla/format"
	"math"
	"reflect"
	"strconv"
	"strings"
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

// DeltaContext represents a context for calculating deltas of scalar values by comparing the current metric values
// against previously known data. It is the callers duty to store previous metric values and provide them to this
// context as a pointer.
type DeltaContext struct {
	*ScalarContext

	previousValue *float64
}

// StringMatchContext represents a context for string values with support for matching against a list of expected values
type StringMatchContext struct {
	*BaseContext

	problemState   State
	expectedValues []string
}

// StringInfoContext represents a context for string values which do not influence the check result but provide
// additional information about the check execution.
type StringInfoContext struct {
	*BaseContext
}

// NewContext instantiates 'BaseContext' with a given name and format string.
func NewContext(name string, format string) *BaseContext {
	return &BaseContext{
		name:          name,
		format:        format,
		resultFactory: NewResultFactory(),
	}
}

// Name represents a getter for the 'name' attribute.
func (c *BaseContext) Name() string {
	return c.name
}

// ResultFactory represents a getter for the 'resultFactory' attribute.
func (c *BaseContext) ResultFactory() ResultFactory {
	return c.resultFactory
}

// SetResultFactory allows overriding the default result factory, which by default instantiates 'BaseResult'.
func (c *BaseContext) SetResultFactory(resultFactory ResultFactory) {
	c.resultFactory = resultFactory
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

// NewDeltaContext instantiates 'DeltaContext' with the given name, pointer for getting/setting previous value and
// optional warning/critical threshold ranges which get applied to the delta values.
func NewDeltaContext(name string, previousValue *float64, warningRange *Range, criticalRange *Range) *DeltaContext {
	return &DeltaContext{
		ScalarContext: NewScalarContext(name, warningRange, criticalRange),

		previousValue: previousValue,
	}
}

// Evaluate calculates the delta between the current value of the given metric against the value stored in the pointer
// which got passed when instantiating this context. The 'new' value of the metric gets automatically stored in this
// pointer, before applying warning/threshold ranges (if available) on the calculated delta value.
func (c *DeltaContext) Evaluate(metric Metric, resource Resource) Result {
	numberMetric, ok := metric.(*NumericMetric)
	if !ok {
		return c.resultFactory(StateUnknown, metric, c, resource,
			fmt.Sprintf("DeltaContext can not process metrics of type [%s]", reflect.TypeOf(metric)))
	}

	metricValue := numberMetric.Value()
	deltaValue := float64(0)
	if !math.IsNaN(metricValue) && !math.IsInf(metricValue, -1) && !math.IsInf(metricValue, 1) {
		*c.previousValue = metricValue
		deltaValue = metricValue - *c.previousValue
	}

	deltaMetric := NewNumericMetric(
		numberMetric.Name()+"_delta",
		deltaValue, "", nil,
		numberMetric.ContextName())

	if c.criticalRange != nil && !c.criticalRange.Match(deltaMetric.Value()) {
		return c.resultFactory(StateCritical, deltaMetric, c, resource, c.criticalRange.ViolationHint())
	} else if c.warningRange != nil && !c.warningRange.Match(deltaMetric.Value()) {
		return c.resultFactory(StateWarning, deltaMetric, c, resource, c.warningRange.ViolationHint())
	} else {
		return c.resultFactory(StateOk, deltaMetric, c, resource, "")
	}
}

// Performance returns the performance data for the delta context without any ranges
func (c *DeltaContext) Performance(metric Metric, resource Resource) *PerfData {
	return &PerfData{
		metric:        metric,
		warningRange:  nil,
		criticalRange: nil,
	}
}

// NewStringMatchContext instantiates 'StringMatchContext' with the given name, a list of expected values and the
// desired result state in case none of the expected values match.
func NewStringMatchContext(name string, rawExpectedValues []string, problemState State) *StringMatchContext {
	var expectedValues []string
	for _, expectedValue := range rawExpectedValues {
		expectedValues = append(expectedValues, strings.ToLower(expectedValue))
	}

	return &StringMatchContext{
		BaseContext: NewContext(name, "%<name>s is %<value>s"),

		problemState:   problemState,
		expectedValues: expectedValues,
	}
}

// Evaluate checks if the given metric matches one of the expected strings, which were passed when instantiating this
// context. In case no expected strings were given, this method will always pass with the state "OK". Please note that
// all checks are done case-insensitive. In case no match was found, the desired 'problemState' (also passed during
// instantiation) gets returned instead.
func (c *StringMatchContext) Evaluate(metric Metric, resource Resource) Result {
	stringMetric, ok := metric.(*StringMetric)
	if !ok {
		return c.resultFactory(StateUnknown, metric, c, resource,
			fmt.Sprintf("StringMatchContext can not process metrics of type [%s]", reflect.TypeOf(metric)))
	}

	if len(c.expectedValues) == 0 {
		return c.resultFactory(StateOk, metric, c, resource, "")
	}

	value := strings.ToLower(stringMetric.Value())
	for _, expectedValue := range c.expectedValues {
		if value == expectedValue {
			return c.resultFactory(StateOk, metric, c, resource, "")
		}
	}

	return c.resultFactory(c.problemState, metric, c, resource,
		fmt.Sprintf("got [%s], expected [%s]", value, strings.Join(c.expectedValues, "],[")))
}

// NewStringInfoContext instantiates 'StringInfoContext' with the given name. This context can be used to pass
// additional information about a check when the user requests verbose plugin output.
func NewStringInfoContext(name string) *StringInfoContext {
	return &StringInfoContext{
		BaseContext: NewContext(name, "%<value>s"),
	}
}

// Evaluate returns a Result object for the given metric and resource. In case of 'StringInfoContext', the check result
// will always return 'StateInfo', which ends up in verbose plugin output unlike 'StateOk'.
func (c *StringInfoContext) Evaluate(metric Metric, resource Resource) Result {
	return c.resultFactory(StateInfo, metric, c, resource, "")
}
