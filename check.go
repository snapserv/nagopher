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
	"errors"
	"fmt"
	"reflect"
	"sort"
)

type Check struct {
	name       string
	resources  []Resource
	contexts   map[string]Context
	results    *ResultCollection
	perfData   []*PerfData
	summarizer Summarizer
}

func NewCheck(name string, summarizer Summarizer) *Check {
	return &Check{
		name:       name,
		contexts:   make(map[string]Context),
		results:    NewResultCollection(),
		summarizer: summarizer,
	}
}

func (c *Check) Run(warnings *WarningCollection) error {
	for _, resource := range c.resources {
		if err := c.evaluateResource(warnings, resource); err != nil {
			c.results.Add(NewResult(StateUnknown, nil, nil, resource, err.Error()))
		}
	}

	sort.Slice(c.perfData, func(i int, j int) bool {
		return c.perfData[i].metric.Name() < c.perfData[j].metric.Name()
	})

	return nil
}

func (c *Check) AttachResources(resources ...Resource) {
	c.resources = append(c.resources, resources...)
}

func (c *Check) AttachContexts(contexts ...Context) {
	for _, context := range contexts {
		c.contexts[context.Name()] = context
	}
}

func (c *Check) GetState() State {
	return c.results.MostSignificantState()
}

func (c *Check) GetSummary() string {
	if c.results.Count() == 0 {
		return c.summarizer.Empty()
	} else {
		if c.GetState() == StateOk {
			return c.summarizer.Ok(c.results)
		} else {
			return c.summarizer.Problem(c.results)
		}
	}
}

func (c *Check) GetVerboseSummary() []string {
	return c.summarizer.Verbose(c.results)
}

func (c *Check) GetPerfData() []*PerfData {
	return c.perfData
}

func (c *Check) evaluateResource(warnings *WarningCollection, resource Resource) error {
	err, metrics := resource.Probe(warnings)
	if err != nil {
		return err
	}
	if len(metrics) == 0 {
		return errors.New(fmt.Sprintf("nagopher: resource [%s] did not return any metrics",
			reflect.TypeOf(resource)))
	}

	for _, metric := range metrics {
		context, ok := c.contexts[metric.ContextName()]
		if !ok {
			return errors.New(fmt.Sprintf("nagopher: missing context with name [%s]", metric.ContextName()))
		}

		result := context.Evaluate(metric, resource)
		if perfData := context.Performance(metric, resource); perfData != nil {
			c.perfData = append(c.perfData, perfData)
		}

		c.results.Add(result)
	}

	return nil
}
