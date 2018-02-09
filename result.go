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

type ResultFactory func(state State, metric Metric, context Context, resource Resource, hint string) Result

type Result interface {
	State() State
	String() string
}

type ResultCollection struct {
	results []Result
}

func NewResultCollection() *ResultCollection {
	return &ResultCollection{}
}

func (rc *ResultCollection) Add(results ...Result) {
	rc.results = append(rc.results, results...)
	rc.sort()
}

func (rc *ResultCollection) Get() []Result {
	return rc.results
}

func (rc *ResultCollection) Count() int {
	return len(rc.results)
}

func (rc *ResultCollection) MostSignificantState() State {
	if len(rc.results) >= 1 {
		return rc.results[0].State()
	} else {
		return StateUnknown
	}
}

func (rc *ResultCollection) sort() {
	sort.Slice(rc.results, func(i, j int) bool {
		return rc.results[i].State().ExitCode > rc.results[j].State().ExitCode
	})
}

type BaseResult struct {
	Result
	state    State
	hint     string
	metric   Metric
	context  Context
	resource Resource
}

func NewResult(state State, metric Metric, context Context, resource Resource, hint string) *BaseResult {
	return &BaseResult{
		state:    state,
		hint:     hint,
		context:  context,
		resource: resource,
		metric:   metric,
	}
}

func NewResultFactory() ResultFactory {
	return func(state State, metric Metric, context Context, resource Resource, hint string) Result {
		return NewResult(state, metric, context, resource, hint)
	}
}

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
	} else {
		return ""
	}
}

func (r *BaseResult) State() State {
	return r.state
}
