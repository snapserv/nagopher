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

type StringMetric interface {
	Metric

	Value() string
}

type stringMetric struct {
	baseMetric
	value string
}

func NewStringMetric(name string, value string, context string) (StringMetric, error) {
	baseMetric, err := newBaseMetric(name, "", nil, context)
	if err != nil {
		return nil, err
	}

	stringMetric := &stringMetric{
		baseMetric: *baseMetric,
		value:      value,
	}

	return stringMetric, nil
}

func (m stringMetric) ToNagiosValue() string {
	return m.ValueString()
}

func (m stringMetric) Value() string {
	return m.value
}

func (m stringMetric) ValueString() string {
	return m.value
}
