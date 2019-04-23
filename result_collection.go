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
	"github.com/markphelps/optional"
	"sort"
)

// ResultCollection contains an arbitrary amount of Result instances and methods to sort them by relevance
type ResultCollection interface {
	Add(results ...Result)
	Get() []Result
	Count() int
	MostSignificantResult() OptionalResult
	MostSignificantState() OptionalState

	GetByMetricName(name string) OptionalResult
	GetMetricByName(name string) OptionalMetric
	GetNumericMetricValue(name string) optional.Float64
	GetStringMetricValue(name string) optional.String
}

type resultCollection struct {
	results []Result
}

// NewResultCollection instantiates a new ResultCollection object without any items
func NewResultCollection() ResultCollection {
	return &resultCollection{}
}

func (c *resultCollection) Add(results ...Result) {
	c.results = append(c.results, results...)
	c.sort()
}

func (c resultCollection) Get() []Result {
	return c.results
}

func (c resultCollection) Count() int {
	return len(c.results)
}

func (c resultCollection) MostSignificantResult() OptionalResult {
	if len(c.results) >= 1 {
		return NewOptionalResult(c.results[0])
	}

	return OptionalResult{}
}

func (c resultCollection) MostSignificantState() OptionalState {
	mostSignificantResult := c.MostSignificantResult()
	if result, err := mostSignificantResult.Get(); err == nil {
		return result.State()
	}

	return OptionalState{}
}

func (c resultCollection) GetByMetricName(name string) OptionalResult {
	for _, result := range c.results {
		metric, err := result.Metric().Get()
		if err != nil || metric == nil {
			continue
		}

		if metric.Name() == name {
			return NewOptionalResult(result)
		}
	}

	return OptionalResult{}
}

func (c resultCollection) GetMetricByName(name string) OptionalMetric {
	result, err := c.GetByMetricName(name).Get()
	if err == nil && result != nil {
		return result.Metric()
	}

	return OptionalMetric{}
}

func (c resultCollection) GetNumericMetricValue(name string) optional.Float64 {
	metric, err := c.GetMetricByName(name).Get()
	if err == nil && metric != nil {
		if numericMetric, ok := metric.(NumericMetric); ok {
			return optional.NewFloat64(numericMetric.Value())
		}
	}

	return optional.Float64{}
}

func (c resultCollection) GetStringMetricValue(name string) optional.String {
	metric, err := c.GetMetricByName(name).Get()
	if err == nil && metric != nil {
		if stringMetric, ok := metric.(StringMetric); ok {
			return optional.NewString(stringMetric.Value())
		}
	}

	return optional.String{}
}

func (c *resultCollection) sort() {
	sort.SliceStable(c.results, func(a int, b int) bool {
		stateA, errA := c.results[a].State().Get()
		stateB, errB := c.results[b].State().Get()

		if errA == nil && errB != nil {
			return true
		} else if errA == nil && errB == nil {
			return stateA.ExitCode() > stateB.ExitCode()
		}

		return false
	})
}
