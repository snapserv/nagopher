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
	"sort"
)

// ResultCollection represents a struct holding 0-n 'Result' objects always guaranteeing to return the collected results
// ordered by the significance of the result states.
type ResultCollection struct {
	results []Result
}

// NewResultCollection instantiates 'ResultCollection', which is by default empty.
func NewResultCollection() *ResultCollection {
	return &ResultCollection{}
}

// Add adds one or more results to the collection and sorts the available results afterwards. In case you want to add a
// lot of results simultaneously, you should only call this method once to avoid unnecessary sorting in between.
func (rc *ResultCollection) Add(results ...Result) {
	rc.results = append(rc.results, results...)
	rc.sort()
}

// Get returns all collected results, being effectively a getter for the 'results' attribute.
func (rc *ResultCollection) Get() []Result {
	return rc.results
}

// GetByMetricName tries to find a result where the associated metric has the given name. Returns nil if no metric could
// be found in the current result collection.
func (rc *ResultCollection) GetByMetricName(name string) Result {
	for _, result := range rc.results {
		metric := result.Metric()
		if metric != nil && metric.Name() == name {
			return result
		}
	}

	return nil
}

// Count is a helper method which returns the current amount of results.
func (rc *ResultCollection) Count() int {
	return len(rc.results)
}

// MostSignificantState returns the most significant state within the result collection. In case no results were added,
// 'StateUnknown' is being returned.
func (rc *ResultCollection) MostSignificantState() State {
	if len(rc.results) >= 1 {
		return rc.results[0].State()
	}

	return StateUnknown
}

func (rc *ResultCollection) sort() {
	sort.Slice(rc.results, func(i, j int) bool {
		return rc.results[i].State().ExitCode > rc.results[j].State().ExitCode
	})
}

// ResultFactory represents a generic function declaration for instantiating 'Result' objects, which can be used to
// override which result type is being used within a 'Context' object.
type ResultFactory func(state State, metric Metric, context Context, resource Resource, hint string) Result

// Result represents a interface for all result types
type Result interface {
	String() string
	State() State
	Metric() Metric
	Hint() string
}

// BaseResult represents a generic context from which all other result types should originate.
type BaseResult struct {
	Result
	state    State
	hint     string
	metric   Metric
	context  Context
	resource Resource
}

// NewResult instantiates 'BaseResult' with the given state, metric, context, resource and hint.
func NewResult(state State, metric Metric, context Context, resource Resource, hint string) *BaseResult {
	return &BaseResult{
		state:    state,
		hint:     hint,
		context:  context,
		resource: resource,
		metric:   metric,
	}
}

// NewResultFactory returns an anonymous function for instantiating a new 'BaseResult'. This is being used by 'Context'
// objects as a generic interface for creating results, which can be backed by various result types.
func NewResultFactory() ResultFactory {
	return func(state State, metric Metric, context Context, resource Resource, hint string) Result {
		return NewResult(state, metric, context, resource, hint)
	}
}

// String returns a readable string for the currently available result values, including both the description (context
// if available, falling back to string representation of metric) and hint if available. In case no values are
// available, an empty string will be returned.
func (r *BaseResult) String() string {
	var description string

	if r.context != nil {
		description = r.context.Describe(r.metric)
	}
	if description == "" && r.metric != nil {
		description = r.metric.String()
	}

	if r.hint != "" && description != "" {
		return fmt.Sprintf("%s (%s)", description, r.hint)
	} else if r.hint != "" {
		return r.hint
	} else if description != "" {
		return description
	}

	return ""
}

// State represents a getter for the 'state' attribute.
func (r *BaseResult) State() State {
	return r.state
}

// Metric represents a getter for the 'metric' attribute.
func (r *BaseResult) Metric() Metric {
	return r.metric
}

// Hint represents a getter for the 'hint' attribute.
func (r *BaseResult) Hint() string {
	return r.hint
}
