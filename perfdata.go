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
	"regexp"
	"strings"
)

type PerfData struct {
	metric        Metric
	warningRange  *Range
	criticalRange *Range
}

func NewPerfData(name string, value float64, valueUnit string, valueRange *Range, warningRange *Range, criticalRange *Range) (error, *PerfData) {
	if strings.ContainsAny(name, "'=") {
		return errors.New(fmt.Sprintf("nagopher: perfdata name [%s] contains invalid characters", name)), nil
	}

	return nil, &PerfData{
		metric:        NewMetric(name, value, valueUnit, valueRange, "perfdata"),
		warningRange:  warningRange,
		criticalRange: criticalRange,
	}
}

func (pd *PerfData) BuildOutput() (error, string) {
	err, quotedName := pd.quoteString(pd.metric.Name())
	if err != nil {
		return err, ""
	}

	output := []string{fmt.Sprintf(
		"%s=%s",
		quotedName,
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

	return nil, strings.TrimRight(strings.Join(output, ";"), ";")
}

func (pd *PerfData) quoteString(value string) (error, string) {
	match, err := regexp.MatchString("^\\w+$", value)
	if err != nil {
		return errors.New(fmt.Sprintf("nagopher: unexpected regexp error (%s)", err.Error())), ""
	}

	if match {
		return nil, value
	} else {
		return nil, fmt.Sprintf("'%s'", value)
	}
}
