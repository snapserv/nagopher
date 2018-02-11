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
	"math"
	"strconv"
	"strings"
)

// Range represents a threshold range with start and end values. Matching can be inverted from inside to outside.
type Range struct {
	invert bool
	start  float64
	end    float64
}

// NewRange instantiates 'Range' with the given start and end values. Additionally, the 'invert' bool defines if the
// 'Match' function should be true when a value is inside (invert = false) or outside (invert = true) the range.
func NewRange(invert bool, start float64, end float64) *Range {
	return &Range{
		invert: invert,
		start:  start,
		end:    end,
	}
}

// ParseRange parses a string representing a Nagios threshold range and instantiates an appropriate Range object.
func ParseRange(specifier string) (*Range, error) {
	var err error
	var start, end = math.NaN(), math.NaN()

	// Invert match if specifier starts with '@' character
	invert := strings.HasPrefix(specifier, "@")
	if invert {
		specifier = specifier[1:]
	}

	// Split specifier by colon
	parts := strings.Split(specifier, ":")
	if len(parts) == 1 {
		if start, err = parseRangePart("", true); err != nil {
			return nil, err
		}
		if end, err = parseRangePart(parts[0], false); err != nil {
			return nil, err
		}
	} else if len(parts) == 2 {
		if start, err = parseRangePart(parts[0], true); err != nil {
			return nil, err
		}
		if end, err = parseRangePart(parts[1], false); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("nagopher: range specifier contains more than one colon")
	}

	// Return new 'Range' instance with parsed values
	return NewRange(invert, start, end), nil
}

// String returns a string representing the current threshold range formatted according to Nagios plugin specifications.
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

// Match returns a boolean if the given value matches the range.
func (r *Range) Match(value float64) bool {
	if value < r.start || value > r.end {
		return r.invert
	}

	return !r.invert
}

// Start returns the start value of the range, formatted according to Nagios plugin specifications. Please note that the
// 'invert' attribute will be disregarded.
func (r *Range) Start() string {
	if r.start == 0 {
		// Zero is automatically being assumed if not given
		return ""
	} else if math.IsInf(r.start, -1) {
		return "~"
	}

	return strconv.FormatFloat(r.start, 'f', -1, strconv.IntSize)
}

// End returns the end value of the range, formatted according to Nagios plugin specifications. Please note that the
// 'invert' attribute will be disregarded.
func (r *Range) End() string {
	if math.IsInf(r.end, 1) {
		return ""
	}

	return strconv.FormatFloat(r.end, 'f', -1, strconv.IntSize)
}

// ViolationHint returns a human-readable string which can be used for describing why a previous value did not match.
// An example output would be: 'outside range 1:10'
func (r *Range) ViolationHint() string {
	start, end := r.Start(), r.End()
	if start == "" {
		start = "0"
	}
	if start == "~" {
		start = "-inf"
	}
	if end == "" {
		end = "inf"
	}

	return fmt.Sprintf("outside range %s:%s", start, end)
}

func parseRangePart(part string, isStart bool) (float64, error) {
	if part == "" {
		if isStart {
			return 0, nil
		}

		return math.Inf(1), nil
	} else if part == "~" {
		if isStart {
			return math.Inf(-1), nil
		}

		return -1, errors.New("nagopher: can not use negative infinity in range as 'end'")
	}

	value, err := strconv.ParseFloat(part, strconv.IntSize)
	if err != nil {
		return -1, fmt.Errorf("nagopher: could not parse range part [%s] as float (%s)", part, err.Error())
	}

	return value, nil
}
