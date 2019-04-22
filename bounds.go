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

//go:generate optional -type=Bounds

import (
	"errors"
	"fmt"
	"github.com/markphelps/optional"
	"math"
	"strconv"
	"strings"
)

// Bounds support matching float64 numbers against a given range, optionally inverting the result
type Bounds interface {
	fmt.Stringer

	ToNagiosRange() string
	Match(value float64) bool

	IsInverted() bool
	Lower() optional.Float64
	Upper() optional.Float64
}

// BoundsOpt is a type alias for functional options used by NewBounds()
type BoundsOpt func(*bounds)

type bounds struct {
	inverted   bool
	lowerBound optional.Float64
	upperBound optional.Float64
}

// NewBounds instantiates Bounds with the given functional options
func NewBounds(options ...BoundsOpt) Bounds {
	threshold := &bounds{
		inverted: false,
	}

	for _, option := range options {
		option(threshold)
	}

	return threshold
}

// FromNagiosRange is a helper method, which constructs a new Bounds object from a Nagios range specifier
func FromNagiosRange(specifier string) (Bounds, error) {
	options, err := NagiosRange(specifier)
	if err != nil {
		return nil, err
	}

	return NewBounds(options...), nil
}

// NagiosRange is a functional option for NewBounds(), which parses a Nagios range specifier
func NagiosRange(specifier string) ([]BoundsOpt, error) {
	var options []BoundsOpt
	var lowerPart, upperPart string

	// Invert bounds if specifier starts with '@'
	if strings.HasPrefix(specifier, "@") {
		options = append(options, InvertedBounds(true))
		specifier = specifier[1:]
	}

	// Split specifier by colon to determine given parts
	parts := strings.Split(specifier, ":")
	if len(parts) == 1 {
		lowerPart, upperPart = parts[0], ""
	} else if len(parts) == 2 {
		lowerPart, upperPart = parts[0], parts[1]
	} else {
		return []BoundsOpt{}, errors.New("range specifier must contain only one colon")
	}

	// Attempt to parse lower bound
	lowerBound, err := parseNagiosRangePart(lowerPart, true)
	if err != nil {
		return []BoundsOpt{}, err
	}
	options = append(options, LowerBound(lowerBound))

	// Attempt to parse upper bound
	upperBound, err := parseNagiosRangePart(upperPart, false)
	if err != nil {
		return []BoundsOpt{}, err
	}
	options = append(options, UpperBound(upperBound))

	// Return bounds options for usage with constructor
	return options, nil
}

func parseNagiosRangePart(rangePart string, isStart bool) (float64, error) {
	if rangePart == "" {
		if isStart {
			return 0, nil
		}

		return math.Inf(1), nil
	}

	if rangePart == "~" {
		if isStart {
			return math.Inf(-1), nil
		}

		return math.NaN(), errors.New("can not use negative infinity in 'end' range")
	}

	value, err := strconv.ParseFloat(rangePart, strconv.IntSize)
	if err != nil {
		return math.NaN(), fmt.Errorf("could not parse range part [%s] as float (%s)", rangePart, err.Error())
	}

	return value, nil
}

// InvertedBounds is a functional option for NewBounds(), which inverts matching of the boundary (inside -> outside)
func InvertedBounds(state bool) BoundsOpt {
	return func(b *bounds) {
		b.inverted = state
	}
}

// LowerBound is a functional option for NewBounds(), which sets the lower boundary
func LowerBound(value float64) BoundsOpt {
	return func(b *bounds) {
		b.lowerBound = optional.NewFloat64(value)
	}
}

// UpperBound is a functional option for NewBounds(), which sets the upper boundary
func UpperBound(value float64) BoundsOpt {
	return func(b *bounds) {
		b.upperBound = optional.NewFloat64(value)
	}
}

func (b bounds) String() string {
	var result string

	lowerBound := b.lowerBound.OrElse(math.NaN())
	upperBound := b.upperBound.OrElse(math.NaN())

	result += strconv.FormatFloat(lowerBound, 'f', -1, strconv.IntSize)
	result += ":"
	result += strconv.FormatFloat(upperBound, 'f', -1, strconv.IntSize)

	if b.inverted {
		return "outside " + result
	}

	return "inside " + result
}

func (b bounds) ToNagiosRange() string {
	var result string

	if b.inverted {
		result += "@"
	}

	if lowerBound, err := b.lowerBound.Get(); err == nil {
		if math.IsInf(lowerBound, -1) {
			result += "~"
		} else if lowerBound != 0 {
			result += strconv.FormatFloat(lowerBound, 'f', -1, strconv.IntSize)
		}
	}

	if upperBound, err := b.upperBound.Get(); err == nil {
		if !math.IsInf(upperBound, 1) {
			result += ":" + strconv.FormatFloat(upperBound, 'f', -1, strconv.IntSize)
		}
	}

	return result
}

func (b bounds) Match(value float64) bool {
	if math.IsNaN(value) || math.IsInf(value, -1) || math.IsInf(value, 1) {
		return false
	}

	lowerBound, err := b.lowerBound.Get()
	if err == nil && value < lowerBound {
		return b.inverted
	}

	upperBound, err := b.upperBound.Get()
	if err == nil && value > upperBound {
		return b.inverted
	}

	return !b.inverted
}

func (b bounds) IsInverted() bool {
	return b.inverted
}

func (b bounds) Lower() optional.Float64 {
	return b.lowerBound
}

func (b bounds) Upper() optional.Float64 {
	return b.upperBound
}
