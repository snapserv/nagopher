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
	// Invert match if specifier starts with '@' character
	invert := strings.HasPrefix(specifier, "@")
	if invert {
		specifier = specifier[1:]
	}

	// Split specifier by colon
	parts := strings.Split(specifier, ":")
	start := parseRangePart(parts[0])
	end := parseRangePart(parts[1])

	// Return new 'Range' instance with parsed values
	return NewRange(invert, start, end)
}

func parseRangePart(part string) float64 {
	if part == "" {
		return math.Inf(1)
	} else if part == "~" {
		return math.Inf(-1)
	} else {
		value, err := strconv.ParseFloat(part, strconv.IntSize)
		if err != nil {
			// TODO: Proper error handling
			panic(fmt.Sprintf("Could not parse part of range specifier: %s (%s)", part, err.Error()))
		}

		return value
	}
}

func (r *Range) String() string {
	return r.Start() + ":" + r.End()
}

func (r *Range) Match(value float64) bool {
	if value < r.start || value > r.end {
		return r.invert
	} else {
		return !r.invert
	}
}

func (r *Range) Start() string {
	if math.IsInf(r.start, -1) {
		return "~"
	} else if math.IsInf(r.start, 1) {
		return ""
	} else {
		return fmt.Sprintf("%g", r.start)
	}
}

func (r *Range) End() string {
	if math.IsInf(r.end, -1) {
		return "~"
	} else if math.IsInf(r.end, 1) {
		return ""
	} else {
		return fmt.Sprintf("%g", r.end)
	}
}

func (r *Range) ViolationHint() string {
	return fmt.Sprintf("outside range %s", r.String())
}
