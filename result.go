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

//go:generate optional -type=Result
package nagopher

import (
	"fmt"
)

type Result interface {
	fmt.Stringer

	Hint() string
	State() OptionalStateData
	Metric() OptionalMetric
	Context() OptionalContext
	Resource() OptionalResource
}

type result struct {
	hint     string
	state    OptionalStateData
	metric   OptionalMetric
	context  OptionalContext
	resource OptionalResource
}
type resultOpt func(*result)

func NewResult(options ...resultOpt) Result {
	result := &result{}

	for _, option := range options {
		option(result)
	}

	return result
}

func ResultHint(value string) resultOpt {
	return func(r *result) {
		r.hint = value
	}
}

func ResultState(value StateData) resultOpt {
	return func(r *result) {
		if value != nil {
			r.state = NewOptionalStateData(value)
		}
	}
}

func ResultMetric(value Metric) resultOpt {
	return func(r *result) {
		if value != nil {
			r.metric = NewOptionalMetric(value)
		}
	}
}

func ResultContext(value Context) resultOpt {
	return func(r *result) {
		if value != nil {
			r.context = NewOptionalContext(value)
		}
	}
}

func ResultResource(value Resource) resultOpt {
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

func (r result) State() OptionalStateData {
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
