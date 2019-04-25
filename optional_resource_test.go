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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptionalResource_Empty(t *testing.T) {
	// when
	optionalResource := OptionalResource{}
	actualResource, err := optionalResource.Get()

	// then
	assert.Error(t, err)
	assert.Empty(t, actualResource)
}

func TestNewOptionalResource(t *testing.T) {
	// given
	expectedResource := NewResource()

	// when
	optionalResource := NewOptionalResource(expectedResource)
	actualResource, err := optionalResource.Get()

	// then
	assert.NoError(t, err)
	assert.Equal(t, true, optionalResource.Present())
	assert.Equal(t, expectedResource, actualResource)
}

func TestOptionalResource_OrElse(t *testing.T) {
	// given
	expectedResource := NewResource()
	alternativeResource := NewResource()

	// when
	optionalResource1 := NewOptionalResource(expectedResource)
	optionalResource2 := OptionalResource{}

	// then
	assert.Equal(t, expectedResource, optionalResource1.OrElse(alternativeResource))
	assert.Equal(t, alternativeResource, optionalResource2.OrElse(alternativeResource))
}

func TestOptionalResource_Set(t *testing.T) {
	// given
	resource1 := NewResource()
	resource2 := NewResource()
	alternativeResource := NewResource()

	// when
	optionalResource := NewOptionalResource(resource1)
	optionalResource.Set(resource2)

	// then
	assert.Equal(t, true, optionalResource.Present())
	assert.Equal(t, resource2, optionalResource.OrElse(alternativeResource))
}

func TestOptionalResource_If(t *testing.T) {
	// given
	var actualResourcePtr *Resource
	var expectedResource Resource
	expectedResource = NewResource()

	// when
	optionalResource := NewOptionalResource(expectedResource)
	optionalResource.If(func(resource Resource) {
		actualResourcePtr = &resource
	})

	// then
	assert.Equal(t, &expectedResource, actualResourcePtr)
}
