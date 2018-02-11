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

import "fmt"

// Summarizer represents a interface for all summary types.
type Summarizer interface {
	Ok(*ResultCollection) string
	Problem(*ResultCollection) string
	Verbose(*ResultCollection) []string
	Empty() string
}

// BaseSummary represents a generic context from which all other result types should originate.
type BaseSummary struct{}

// NewBaseSummary instantiates 'BaseSummary'.
func NewBaseSummary() *BaseSummary {
	return &BaseSummary{}
}

// Ok returns the string representation of the most-significant result, which is always the first element in a
// 'ResultCollection' object. This method should be called if the global check state equals to 'StateOk'.
func (s *BaseSummary) Ok(resultCollection *ResultCollection) string {
	return resultCollection.Get()[0].String()
}

// Problem returns the string representation of the most-significant result, which is always the first element in a
// 'ResultCollection' object. This method should be called if the global check state does not equal to 'StateOk'.
func (s *BaseSummary) Problem(resultCollection *ResultCollection) string {
	return resultCollection.Get()[0].String()
}

// Verbose returns the string representation of all results which do not have a state equal to 'StateOk'.
func (s *BaseSummary) Verbose(resultCollection *ResultCollection) []string {
	var messages []string

	for _, result := range resultCollection.Get() {
		if result.State() == StateOk {
			continue
		}

		messages = append(messages, fmt.Sprintf("%s: %s", result.State().Description, result))
	}

	return messages
}

// Empty returns the string representation when no results are available. This method should be called when the caller
// does not have any result collection or the result collection is empty.
func (s *BaseSummary) Empty() string {
	return "No check results"
}
