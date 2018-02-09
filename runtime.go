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
	"os"
	"strings"
)

type Runtime struct {
	fmt.Stringer
	verbose bool
}

type CheckResult struct {
	ExitCode int
	Output   string
}

var illegalCharacters = []string{"|"}

func NewRuntime(verbose bool) *Runtime {
	return &Runtime{
		verbose: verbose,
	}
}

func (o *Runtime) Execute(check *Check) CheckResult {
	check.Run()

	return CheckResult{
		ExitCode: check.GetState().ExitCode,
		Output:   o.buildCheckOutput(check),
	}
}

func (o *Runtime) ExecuteAndExit(check *Check) {
	result := o.Execute(check)
	fmt.Print(result.Output)
	os.Exit(result.ExitCode)
}

func (o *Runtime) buildCheckOutput(check *Check) string {
	var perfData string

	output := o.buildCheckStatusOutput(check)
	if o.verbose {
		output += "\n" + strings.Join(o.filterStrings(check.GetVerboseSummary()), "\n")
		perfData = o.buildCheckPerfDataOutput(check.GetPerfData(), "\n")
	} else {
		perfData = o.buildCheckPerfDataOutput(check.GetPerfData(), " ")
	}

	if perfData != "" {
		output += " | " + perfData
	}

	return output
}

func (o *Runtime) buildCheckStatusOutput(check *Check) string {
	var output []string

	if check.name != "" {
		output = append(output, strings.ToUpper(check.name))
	}
	output = append(output, strings.ToUpper(check.GetState().Description))
	if summary := strings.TrimSpace(check.GetSummary()); summary != "" {
		output = append(output, "-", summary)
	}

	return o.filterString(strings.Join(output, " "))
}

func (o *Runtime) buildCheckPerfDataOutput(perfData []*PerfData, separator string) string {
	perfDataStrings := make([]string, len(perfData))
	for key, value := range perfData {
		perfDataStrings[key] = value.String()
	}

	return o.filterString(strings.Join(perfDataStrings, separator))
}

func (o *Runtime) filterString(value string) string {
	for _, character := range illegalCharacters {
		value = strings.Replace(value, character, "", -1)
	}

	// TODO: Print warning for filtered characters

	return value
}

func (o *Runtime) filterStrings(values []string) []string {
	results := make([]string, len(values))
	for key, value := range values {
		results[key] = o.filterString(value)
	}

	return results
}
