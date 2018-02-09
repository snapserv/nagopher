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
	"math"
	"strconv"
	"strings"
)

type Range struct {
	invert bool
	start  float64
	end    float64
}

func NewRange(invert bool, start float64, end float64) *Range {
	return &Range{
		invert: invert,
		start:  start,
		end:    end,
	}
}

func ParseRange(specifier string) *Range {
	var start float64 = math.NaN()
	var end float64 = math.NaN()

	// Invert match if specifier starts with '@' character
	invert := strings.HasPrefix(specifier, "@")
	if invert {
		specifier = specifier[1:]
	}

	// Split specifier by colon
	parts := strings.Split(specifier, ":")
	if len(parts) == 1 {
		end = parseRangePart(parts[0])
	} else if len(parts) == 2 {
		start = parseRangePart(parts[0])
		end = parseRangePart(parts[1])
	} else {
		// TODO: Proper error handling without panicking
		panic("Invalid amount of colons...")
	}

	// Handle missing values by setting them to default values
	if math.IsNaN(start) {
		start = 0
	}
	if math.IsNaN(end) {
		end = math.Inf(1)
	}

	// Return new 'Range' instance with parsed values
	return NewRange(invert, start, end)
}

func parseRangePart(part string) float64 {
	if part == "" {
		return math.NaN()
	} else if part == "~" {
		return math.Inf(-1)
	} else {
		value, err := strconv.ParseFloat(part, strconv.IntSize)
		if err != nil {
			// TODO: Proper error handling without panicking
			panic(fmt.Sprintf("Could not parse part of range specifier: %s (%s)", part, err.Error()))
		}

		return value
	}
}

func (r *Range) String() string {
	invertPrefix := ""
	start, end := r.Start(), r.End()

	if r.invert {
		invertPrefix = "@"
	}

	if start == "" && end == "" {
		return ""
	} else if start == "" && end != "" {
		return invertPrefix + end
	} else {
		return invertPrefix + start + ":" + end
	}
}

func (r *Range) Match(value float64) bool {
	if value < r.start || value > r.end {
		return r.invert
	} else {
		return !r.invert
	}
}

func (r *Range) Start() string {
	if r.start == 0 {
		// Zero is automatically being assumed if not given
		return ""
	} else if math.IsInf(r.start, -1) {
		return "~"
	} else if math.IsInf(r.start, 1) {
		return ""
	} else {
		return strconv.FormatFloat(r.start, 'f', -1, strconv.IntSize)
	}
}

func (r *Range) End() string {
	if math.IsInf(r.end, -1) {
		return "~"
	} else if math.IsInf(r.end, 1) {
		return ""
	} else {
		return strconv.FormatFloat(r.end, 'f', -1, strconv.IntSize)
	}
}

func (r *Range) ViolationHint() string {
	return fmt.Sprintf("outside range %s", r.String())
}
