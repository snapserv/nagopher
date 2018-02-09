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
	"regexp"
	"strings"
)

type PerfData struct {
	metric        Metric
	warningRange  *Range
	criticalRange *Range
}

func NewPerfData(name string, value float64, valueUnit string, valueRange *Range, warningRange *Range, criticalRange *Range) *PerfData {
	if strings.ContainsAny(name, "'=") {
		// TODO: Illegal characters, abort here
	}

	return &PerfData{
		metric:        NewMetric(name, value, valueUnit, valueRange, "perfdata"),
		warningRange:  warningRange,
		criticalRange: criticalRange,
	}
}

func (pd *PerfData) String() string {
	output := []string{fmt.Sprintf(
		"%s=%s",
		pd.quoteString(pd.metric.Name()),
		pd.metric.ValueUnit(),
	)}

	if pd.warningRange != nil {
		output = append(output, pd.warningRange.String())
	} else {
		output = append(output, "")
	}

	if pd.criticalRange != nil {
		output = append(output, pd.criticalRange.String())
	} else {
		output = append(output, "")
	}

	if valueRange := pd.metric.ValueRange(); valueRange != nil {
		if start := valueRange.Start(); start != "" {
			output = append(output, start)
		} else {
			output = append(output, "")
		}

		if end := valueRange.End(); end != "" {
			output = append(output, end)
		} else {
			output = append(output, "")
		}
	} else {
		output = append(output, "", "")
	}

	return strings.TrimRight(strings.Join(output, ";"), ";")
}

func (pd *PerfData) quoteString(value string) string {
	match, err := regexp.MatchString("^\\w+$", value)
	if err != nil {
		panic("Unexpected runtime error: " + err.Error())
	}

	if match {
		return value
	} else {
		return fmt.Sprintf("'%s'", value)
	}
}
