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

type Summarizer interface {
	Ok(*ResultCollection) string
	Problem(*ResultCollection) string
	Verbose(*ResultCollection) []string
	Empty() string
}

type Summary struct{}

func NewSummary() *Summary {
	return &Summary{}
}

func (s *Summary) Ok(resultCollection *ResultCollection) string {
	return resultCollection.Get()[0].String()
}

func (s *Summary) Problem(resultCollection *ResultCollection) string {
	return resultCollection.Get()[0].String()
}

func (s *Summary) Verbose(resultCollection *ResultCollection) []string {
	var messages []string

	for _, result := range resultCollection.Get() {
		if result.State() == StateOk {
			continue
		}

		messages = append(messages, fmt.Sprintf("%s: %s", result.State().Description, result))
	}

	return messages
}

func (s *Summary) Empty() string {
	return "No check results"
}
