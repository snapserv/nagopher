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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseMetric_ContextName(t *testing.T) {
	metric := NewBaseMetric("metric", "", nil, "context")
	assert.Equal(t, "context", metric.ContextName())
}

func TestNumberMetric_ValueUnit_Float(t *testing.T) {
	metric := NewNumericMetric("time", 3.141, "s", nil, "")
	assert.Equal(t, "3.141s", metric.ValueUnit())
}

func TestNumberMetric_ValueUnit_Integer(t *testing.T) {
	metric := NewNumericMetric("count", 42, "", nil, "")
	assert.Equal(t, "42", metric.ValueUnit())
}

func TestNumberMetric_ValueUnit_LargeInteger(t *testing.T) {
	metric := NewNumericMetric("grains", 4200000000, "", nil, "")
	assert.Equal(t, "4200000000", metric.ValueUnit())
}

func TestNumberMetric_ValueUnit_LargeFloat(t *testing.T) {
	metric := NewNumericMetric("grains", 420000.42, "", nil, "")
	assert.Equal(t, "420000.42", metric.ValueUnit())
}

func TestStringMetric_ValueUnit(t *testing.T) {
	metric := NewStringMetric("state", "OK", "")
	assert.Equal(t, "OK", metric.ValueUnit())
}
