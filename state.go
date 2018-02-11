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

// State represents a result state, containing an exit code and human-readable description.
type State struct {
	ExitCode    int
	Description string
}

// StateOk represents an "OK" state with exit code 0.
var StateOk = State{ExitCode: 0, Description: "ok"}

// StateWarning represents a "WARNING" state with exit code 1.
var StateWarning = State{ExitCode: 1, Description: "warning"}

// StateCritical represents a "CRITICAL" state with exit code 2.
var StateCritical = State{ExitCode: 2, Description: "critical"}

// StateUnknown represents an "UNKNOWN" state with exit code 3.
var StateUnknown = State{ExitCode: 3, Description: "unknown"}
