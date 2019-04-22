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

import "sort"

type Collection interface {
	Add(results ...Result)
	Get() []Result
	Count() int
	MostSignificantResult() OptionalResult
	MostSignificantState() OptionalStateData
}

type collection struct {
	results []Result
}

func NewResultCollection() Collection {
	return &collection{}
}

func (c *collection) Add(results ...Result) {
	c.results = append(c.results, results...)
	c.sort()
}

func (c collection) Get() []Result {
	return c.results
}

func (c collection) Count() int {
	return len(c.results)
}

func (c collection) MostSignificantResult() OptionalResult {
	if len(c.results) >= 1 {
		return NewOptionalResult(c.results[0])
	}

	return OptionalResult{}
}

func (c collection) MostSignificantState() OptionalStateData {
	mostSignificantResult := c.MostSignificantResult()
	if result, err := mostSignificantResult.Get(); err == nil {
		return result.State()
	}

	return OptionalStateData{}
}

func (c *collection) sort() {
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
