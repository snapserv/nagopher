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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetric_ValueUnit_Float(t *testing.T) {
	metric := NewMetric("time", 3.141, "s", nil, "")
	assert.Equal(t, "3.141s", metric.ValueUnit())
}

func TestMetric_ValueUnit_Integer(t *testing.T) {
	metric := NewMetric("count", 42, "", nil, "")
	assert.Equal(t, "42", metric.ValueUnit())
}

func TestMetric_ValueUnit_LargeInteger(t *testing.T) {
	metric := NewMetric("grains", 4200000000, "", nil, "")
	assert.Equal(t, "4200000000", metric.ValueUnit())
}

func TestMetric_ValueUnit_LargeFloat(t *testing.T) {
	metric := NewMetric("grains", 420000.42, "", nil, "")
	assert.Equal(t, "420000.42", metric.ValueUnit())
}
