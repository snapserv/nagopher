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

func TestCheck_AttachContexts(t *testing.T) {
	context1 := NewContext("ctx1", "")
	context2 := NewContext("ctx2", "")
	check := NewCheck("check", NewBaseSummary())

	check.AttachContexts(context1, context2)
	assert.Equal(t, check.contexts["ctx1"], context1)
	assert.Equal(t, check.contexts["ctx2"], context2)
}

func TestCheck_AttachResources(t *testing.T) {
	resource1 := NewResource()
	resource2 := NewResource()
	check := NewCheck("check", NewBaseSummary())

	check.AttachResources(resource1, resource2)
	assert.Contains(t, check.resources, resource1)
	assert.Contains(t, check.resources, resource2)
}
