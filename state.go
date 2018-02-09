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

type State struct {
	ExitCode    int
	Description string
}

var (
	StateOk       = State{ExitCode: 0, Description: "ok"}
	StateWarning  = State{ExitCode: 1, Description: "warning"}
	StateCritical = State{ExitCode: 2, Description: "critical"}
	StateUnknown  = State{ExitCode: 3, Description: "unknown"}
)
