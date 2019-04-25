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

// State represents a Nagios plugin state, which consists of an exit code and description
type State interface {
	ExitCode() int8
	Description() string
}

type state struct {
	exitCode    int8
	description string
}

// StateUnknown returns an "UNKNOWN" state according to Nagios plugin standards
func StateUnknown() State {
	return state{exitCode: 3, description: "unknown"}
}

// StateCritical returns an "CRITICAL" state according to Nagios plugin standards
func StateCritical() State {
	return state{exitCode: 2, description: "critical"}
}

// StateWarning returns an "WARNING" state according to Nagios plugin standards
func StateWarning() State {
	return state{exitCode: 1, description: "warning"}
}

// StateOk returns an "OK" state according to Nagios plugin standards
func StateOk() State {
	return state{exitCode: 0, description: "ok"}
}

func (s state) ExitCode() int8 {
	return s.exitCode
}

func (s state) Description() string {
	return s.description
}
