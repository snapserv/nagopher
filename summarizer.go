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

type Summarizer interface {
	Ok(Check) string
	Problem(Check) string
	Verbose(Check) []string
	Empty() string
}

type baseSummarizer struct{}

func NewSummarizer() *baseSummarizer {
	return &baseSummarizer{}
}

func (s baseSummarizer) Ok(check Check) string {
	result, err := check.Results().MostSignificantResult().Get()
	if err != nil {
		return s.Empty()
	}

	return result.String()
}

func (s baseSummarizer) Problem(check Check) string {
	result, err := check.Results().MostSignificantResult().Get()
	if err != nil {
		return s.Empty()
	}

	return result.String()
}

func (s baseSummarizer) Verbose(check Check) []string {
	var messages []string

	for _, result := range check.Results().Get() {
		state, err := result.State().Get()
		if err != nil || state == nil {
			messages = append(messages, fmt.Sprintf("info: %s", result))
			continue
		}

		if state == StateOk() {
			continue
		}
		messages = append(messages, fmt.Sprintf("%s: %s", state.Description(), result))
	}

	return messages
}

func (s baseSummarizer) Empty() string {
	return "No check results"
}
