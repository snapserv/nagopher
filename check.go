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
	"sort"
)

// Check represents a nagopher check containing all required objects for execution, evaluation and visualization.
type Check struct {
	name       string
	meta       map[string]interface{}
	resources  []Resource
	contexts   map[string]Context
	results    *ResultCollection
	perfData   []*PerfData
	summarizer Summarizer
}

// NewCheck instantiates 'Check' with a given name and summarizer.
func NewCheck(name string, summarizer Summarizer) *Check {
	return &Check{
		name:       name,
		meta:       make(map[string]interface{}),
		contexts:   make(map[string]Context),
		results:    NewResultCollection(),
		summarizer: summarizer,
	}
}

// Run executes all probes of the attached resources and collects their results including performance data.
func (c *Check) Run(warnings *WarningCollection) {
	for _, resource := range c.resources {
		if err := c.evaluateResource(warnings, resource); err != nil {
			c.results.Add(NewResult(StateUnknown, nil, nil, resource, err.Error()))
		}
	}

	sort.Slice(c.perfData, func(i int, j int) bool {
		return c.perfData[i].metric.Name() < c.perfData[j].metric.Name()
	})
}

// SetMeta sets the metadata for a given key to the passed value. These metadata methods can store any arbitrary type by
// accepting 'interface{}' as their type. Please note that this methods should never be used directly within this API,
// as the user is able to override any field without any further checks.
func (c *Check) SetMeta(key string, value interface{}) {
	c.meta[key] = value
}

// GetMeta returns the metadata for a given key or the passed default value in case the key does not exist. Please also
// note the documentation about 'SetMeta' for further information about the metadata storage system.
func (c *Check) GetMeta(key string, defaultValue interface{}) interface{} {
	if value, ok := c.meta[key]; ok {
		return value
	}

	return defaultValue
}

// AttachResources attaches one or more resources to the check.
func (c *Check) AttachResources(resources ...Resource) {
	c.resources = append(c.resources, resources...)
}

// AttachContexts attaches one or more contexts to the check.
func (c *Check) AttachContexts(contexts ...Context) {
	for _, context := range contexts {
		c.contexts[context.Name()] = context
	}
}

// Results represents a getter for the 'results' attribute.
func (c *Check) Results() *ResultCollection {
	return c.results
}

// GetState returns the most significant state based on current check results.
func (c *Check) GetState() State {
	return c.results.MostSignificantState()
}

// GetSummary returns the output of the checks summarizer.
func (c *Check) GetSummary() string {
	if c.results.Count() == 0 {
		return c.summarizer.Empty()
	}

	if c.GetState() == StateOk {
		return c.summarizer.Ok(c)
	}

	return c.summarizer.Problem(c)
}

// GetVerboseSummary returns the verbose output of the checks summarizer.
func (c *Check) GetVerboseSummary() []string {
	return c.summarizer.Verbose(c)
}

// GetPerfData returns the currently available performance data.
func (c *Check) GetPerfData() []*PerfData {
	return c.perfData
}

func (c *Check) evaluateResource(warnings *WarningCollection, resource Resource) error {
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
		if perfData := context.Performance(metric, resource); perfData != nil {
			c.perfData = append(c.perfData, perfData)
		}

		c.results.Add(result)
	}

	return nil
}
