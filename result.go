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

import (
	"fmt"
)

// Result represents the output of a context, after the evaluation of an associated metric
type Result interface {
	fmt.Stringer

	Hint() string
	State() OptionalState
	Metric() OptionalMetric
	Context() OptionalContext
	Resource() OptionalResource
}

// ResultOpt is a type alias for functional options used by NewResult()
type ResultOpt func(*result)

type result struct {
	hint     string
	state    OptionalState
	metric   OptionalMetric
	context  OptionalContext
	resource OptionalResource
}

// NewResult instantiates a new Result with the given functional options
func NewResult(options ...ResultOpt) Result {
	result := &result{}

	for _, option := range options {
		option(result)
	}

	return result
}

// ResultHint is a functional option for NewResult(), which stores the hint of the result
func ResultHint(value string) ResultOpt {
	return func(r *result) {
		r.hint = value
	}
}

// ResultState is a functional option for NewResult(), which stores the state of the result
func ResultState(value State) ResultOpt {
	return func(r *result) {
		if value != nil {
			r.state = NewOptionalState(value)
		}
	}
}

// ResultMetric is a functional option for NewResult(), which stores the responsible metric of the result
func ResultMetric(value Metric) ResultOpt {
	return func(r *result) {
		if value != nil {
			r.metric = NewOptionalMetric(value)
		}
	}
}

// ResultContext is a functional option for NewResult(), which stores the responsible context of the result
func ResultContext(value Context) ResultOpt {
	return func(r *result) {
		if value != nil {
			r.context = NewOptionalContext(value)
		}
	}
}

// ResultResource is a functional option for NewResult(), which stores the responsible resource of the result
func ResultResource(value Resource) ResultOpt {
	return func(r *result) {
		if value != nil {
			r.resource = NewOptionalResource(value)
		}
	}
}

func (r result) String() string {
	var description string

	if metric, err := r.metric.Get(); err == nil {
		if context, err := r.context.Get(); err == nil {
			description = context.Describe(metric)
		} else {
			description = metric.ToNagiosValue()
		}
	}

	if r.hint != "" && description != "" {
		return fmt.Sprintf("%s (%s)", description, r.hint)
	} else if r.hint != "" {
		return r.hint
	}

	return description
}

func (r result) Hint() string {
	return r.hint
}

func (r result) State() OptionalState {
	return r.state
}

func (r result) Metric() OptionalMetric {
	return r.metric
}

func (r result) Context() OptionalContext {
	return r.context
}

func (r result) Resource() OptionalResource {
	return r.resource
}
