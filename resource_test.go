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
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestBaseResourceMockResourceWarnings struct {
	*BaseResource
	warnings []Warning
}

func NewTestBaseResourceMockResourceWarnings(warnings []Warning) *TestBaseResourceMockResourceWarnings {
	return &TestBaseResourceMockResourceWarnings{
		BaseResource: NewResource(),
		warnings:     warnings,
	}
}

func (r *TestBaseResourceMockResourceWarnings) Probe(warnings *WarningCollection) (error, []Metric) {
	warnings.Add(r.warnings...)
	return nil, []Metric{}
}

type TestBaseResourceMockResourceError struct {
	*BaseResource
}

func NewTestBaseResourceMockResourceError() *TestBaseResourceMockResourceError {
	return &TestBaseResourceMockResourceError{
		BaseResource: NewResource(),
	}
}

func (r *TestBaseResourceMockResourceError) Probe(warnings *WarningCollection) (error, []Metric) {
	return errors.New("nagopher-test: probing has failed"), []Metric{}
}

func TestBaseResource_Probe_NoMetrics(t *testing.T) {
	warnings := NewWarningCollection()
	resource := NewResource()
	err, metrics := resource.Probe(warnings)

	assert.Nil(t, err)
	assert.Equal(t, []Metric{}, metrics)
}

func TestBaseResource_Probe_Warnings(t *testing.T) {
	warningCollection := NewWarningCollection()
	warnings := []Warning{
		NewWarning("nagopher-test: first warning"),
		NewWarning("nagopher-test: second warning"),
	}

	resource := NewTestBaseResourceMockResourceWarnings(warnings)
	err, metrics := resource.Probe(warningCollection)

	assert.Nil(t, err)
	assert.Equal(t, []Metric{}, metrics)
	for _, warning := range warnings {
		assert.Contains(t, warningCollection.GetStrings(), warning.Warning())
	}
}

func TestBaseResource_Probe_Error(t *testing.T) {
	warnings := NewWarningCollection()
	resource := NewTestBaseResourceMockResourceError()
	err, metrics := resource.Probe(warnings)

	assert.NotNil(t, err)
	assert.Empty(t, metrics)
}
