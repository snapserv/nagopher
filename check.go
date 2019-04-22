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
	"reflect"
	"sort"
)

type Check interface {
	Run(warnings WarningCollection)
	SetMeta(key string, value interface{})
	GetMeta(key string, defaultValue interface{}) interface{}
	AttachResources(resources ...Resource)
	AttachContexts(contexts ...Context)

	Name() string
	PerfData() []PerfData
	Results() Collection
	State() StateData
	Summary() string
	VerboseSummary() []string
}

type baseCheck struct {
	name         string
	meta         map[string]interface{}
	contexts     map[string]Context
	resources    []Resource
	performances []PerfData
	results      Collection
	summarizer   Summarizer
}

func NewCheck(name string, summarizer Summarizer) Check {
	check := &baseCheck{
		name:       name,
		summarizer: summarizer,
		meta:       make(map[string]interface{}),
		contexts:   make(map[string]Context),
		results:    NewResultCollection(),
	}

	return check
}

func (c *baseCheck) Run(warnings WarningCollection) {
	for _, resource := range c.resources {
		err := c.evaluateResource(warnings, resource)
		if err != nil {
			c.results.Add(NewResult(
				ResultState(StateUnknown()),
				ResultResource(resource), ResultHint(err.Error()),
			))
		}
	}

	sort.SliceStable(c.performances, func(a int, b int) bool {
		return c.performances[a].Metric().Name() < c.performances[b].Metric().Name()
	})
}

func (c *baseCheck) evaluateResource(warnings WarningCollection, resource Resource) error {
	metrics, err := resource.Probe(warnings)
	if err != nil {
		return err
	}
	if len(metrics) == 0 {
		return fmt.Errorf("nagopher: resource [%s] did not return any metrics", reflect.TypeOf(resource))
	}

	for _, metric := range metrics {
		context, ok := c.contexts[metric.ContextName()]
		if !ok {
			return fmt.Errorf("nagopher: missing context with name [%s]", metric.ContextName())
		}

		result := context.Evaluate(metric, resource)
		c.results.Add(result)

		perfData, err := context.Performance(metric, resource)
		if err != nil {
			return fmt.Errorf("nagopher: collecting performance data failed with [%s]", err.Error())
		}
		if performance, err := perfData.Get(); err == nil {
			c.performances = append(c.performances, performance)
		}
	}

	return nil
}

func (c *baseCheck) SetMeta(key string, value interface{}) {
	c.meta[key] = value
}

func (c baseCheck) GetMeta(key string, defaultValue interface{}) interface{} {
	if value, ok := c.meta[key]; ok {
		return value
	}

	return defaultValue
}

func (c *baseCheck) AttachResources(resources ...Resource) {
	c.resources = append(c.resources, resources...)
}

func (c *baseCheck) AttachContexts(contexts ...Context) {
	for _, context := range contexts {
		c.contexts[context.Name()] = context
	}
}

func (c baseCheck) Results() Collection {
	return c.results
}

func (c baseCheck) State() StateData {
	state, err := c.results.MostSignificantState().Get()
	if err == nil {
		return state
	}

	return StateUnknown()
}

func (c baseCheck) Summary() string {
	if c.results.Count() == 0 {
		return c.summarizer.Empty()
	}

	if c.State() == StateOk() {
		return c.summarizer.Ok(&c)
	}

	return c.summarizer.Problem(&c)
}

func (c baseCheck) VerboseSummary() []string {
	return c.summarizer.Verbose(&c)
}

func (c baseCheck) Name() string {
	return c.name
}

func (c baseCheck) PerfData() []PerfData {
	return c.performances
}
