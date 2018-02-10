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

func ParseRange(specifier string) (error, *Range) {
	var err error
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
		if err, start = parseRangePart("", true); err != nil {
			return err, nil
		}
		if err, end = parseRangePart(parts[0], false); err != nil {
			return err, nil
		}
	} else if len(parts) == 2 {
		if err, start = parseRangePart(parts[0], true); err != nil {
			return err, nil
		}
		if err, end = parseRangePart(parts[1], false); err != nil {
			return err, nil
		}
	} else {
		return errors.New("nagopher: range specifier contains more than one colon"), nil
	}

	// Return new 'Range' instance with parsed values
	return nil, NewRange(invert, start, end)
}

func parseRangePart(part string, isStart bool) (error, float64) {
	if part == "" {
		if isStart {
			return nil, 0
		} else {
			return nil, math.Inf(1)
		}
	} else if part == "~" {
		if isStart {
			return nil, math.Inf(-1)
		} else {
			return errors.New("nagopher: can not use negative infinity in range as 'end'"), -1
		}
	} else {
		value, err := strconv.ParseFloat(part, strconv.IntSize)
		if err != nil {
			return errors.New(fmt.Sprintf("nagopher: could not parse range part [%s] as float (%s)",
				part, err.Error())), -1
		}

		return nil, value
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
	} else {
		return strconv.FormatFloat(r.start, 'f', -1, strconv.IntSize)
	}
}

func (r *Range) End() string {
	if math.IsInf(r.end, 1) {
		return ""
	} else {
		return strconv.FormatFloat(r.end, 'f', -1, strconv.IntSize)
	}
}

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
