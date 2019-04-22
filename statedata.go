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

//go:generate optional -type=StateData
package nagopher

type StateData interface {
	ExitCode() int8
	Description() string
}

type stateData struct {
	exitCode    int8
	description string
}

func StateUnknown() StateData {
	return stateData{exitCode: 3, description: "unknown"}
}

func StateCritical() StateData {
	return stateData{exitCode: 2, description: "critical"}
}

func StateWarning() StateData {
	return stateData{exitCode: 1, description: "warning"}
}

func StateOk() StateData {
	return stateData{exitCode: 0, description: "ok"}
}

func (s stateData) ExitCode() int8 {
	return s.exitCode
}

func (s stateData) Description() string {
	return s.description
}
