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
	"fmt"
	"regexp"
	"strings"
)

// PerfData holds a Metric instance and allows transformation into Nagios performance data
type PerfData interface {
	ToNagiosPerfData() string
	Metric() Metric
}

type perfData struct {
	metric            Metric
	warningThreshold  OptionalBounds
	criticalThreshold OptionalBounds
}

const illegalNameChars = "'="

// NewPerfData instantiates a new PerfData with the given metric and optional thresholds
func NewPerfData(metric Metric, warningThreshold *Bounds, criticalThreshold *Bounds) (PerfData, error) {
	if strings.ContainsAny(metric.Name(), illegalNameChars) {
		return nil, fmt.Errorf("perfdata metric name [%s] contains invalid characters", metric.Name())
	}

	perfData := &perfData{
		metric: metric,
	}

	if warningThreshold != nil {
		perfData.warningThreshold = NewOptionalBounds(*warningThreshold)
	}
	if criticalThreshold != nil {
		perfData.criticalThreshold = NewOptionalBounds(*criticalThreshold)
	}

	return perfData, nil
}

// NewNumericPerfData instantiates new PerfData. The parameters are being used to create a new NumericMetric instance,
// which is then passed along with the optional thresholds to NewPerfData().
func NewNumericPerfData(name string, value float64, valueUnit string, valueRange *Bounds,
	warningThreshold *Bounds, criticalThreshold *Bounds) (PerfData, error) {
	numericMetric, err := NewNumericMetric(name, value, valueUnit, valueRange, "perfdata")
	if err != nil {
		return nil, err
	}

	return NewPerfData(numericMetric, warningThreshold, criticalThreshold)
}

func (pd perfData) ToNagiosPerfData() string {
	quotedName := pd.quoteString(pd.metric.Name())
	emptyBounds := NewBounds()

	valueRange := pd.metric.ValueRange().OrElse(emptyBounds)
	valueRangeParts := strings.Split(valueRange.ToNagiosRange(), ":")
	warningThreshold := pd.warningThreshold.OrElse(emptyBounds)
	criticalThreshold := pd.criticalThreshold.OrElse(emptyBounds)

	outputValues := append([]string{
		fmt.Sprintf("%s=%s", quotedName, pd.metric.ToNagiosValue()),
		warningThreshold.ToNagiosRange(),
		criticalThreshold.ToNagiosRange(),
	}, valueRangeParts...)

	output := strings.TrimRight(strings.Join(outputValues, ";"), ";")
	return output
}

func (pd perfData) Metric() Metric {
	return pd.metric
}

func (pd perfData) quoteString(value string) string {
	match := regexp.MustCompile("^\\w+$").MatchString(value)
	if match {
		return value
	}

	return fmt.Sprintf("'%s'", value)
}
