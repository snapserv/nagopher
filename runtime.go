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
	"os"
	"strings"
)

// Runtime represents a framework for executing and outputting checks.
type Runtime struct {
	fmt.Stringer
	verbose bool
}

// CheckResult represents the result of an executed check, represented by an exit code and the string representation of
// the check output according to the Nagios plugin specifications.
type CheckResult struct {
	ExitCode int
	Output   string
}

var illegalCharacters = []string{"|"}

// NewRuntime instantiates 'Runtime' and specifies if the check output should be verbose or not.
func NewRuntime(verbose bool) *Runtime {
	return &Runtime{
		verbose: verbose,
	}
}

// Execute executes a single check and returns a 'CheckResult' object. Any errors occurring within the checks are not
// being passed to the caller, as they will be represented as 'StateUnknown' results. In case an unexpected runtime
// error occurs, this method will call panic - however this should never happen under normal circumstances.
func (o *Runtime) Execute(check *Check) CheckResult {
	warnings := NewWarningCollection()
	check.Run(warnings)

	checkOutput, err := o.buildCheckOutput(check, warnings)
	if err != nil {
		panic(fmt.Sprintf("nagopher: unexpected runtime error [%s]", err.Error()))
	}

	return CheckResult{
		ExitCode: check.GetState().ExitCode,
		Output:   checkOutput,
	}
}

// ExecuteAndExit is a helper method for calling 'Execute', followed by printing the returned check results and exiting
// the current process with the returned exit code.
func (o *Runtime) ExecuteAndExit(check *Check) {
	result := o.Execute(check)
	fmt.Print(result.Output)
	os.Exit(result.ExitCode)
}

func (o *Runtime) buildCheckOutput(check *Check, warnings *WarningCollection) (string, error) {
	output := o.buildCheckStatusOutput(check, warnings)

	if perfData, err := o.buildCheckPerfDataOutput(check.GetPerfData(), " ", warnings); err == nil {
		output += " | " + perfData
	} else {
		return "", err
	}
	output += "\n"

	if o.verbose {
		lines := o.filterStrings(warnings, check.GetVerboseSummary())
		if len(lines) > 0 {
			output += strings.Join(lines, "\n") + "\n"
		}
	}

	if warningStrings := warnings.GetStrings(); len(warningStrings) > 0 {
		output += strings.Join(o.filterStrings(nil, warningStrings), "\n") + "\n"
	}

	return output, nil
}

func (o *Runtime) buildCheckStatusOutput(check *Check, warnings *WarningCollection) string {
	var output []string

	if check.name != "" {
		output = append(output, strings.ToUpper(check.name))
	}
	output = append(output, strings.ToUpper(check.GetState().Description))
	if summary := strings.TrimSpace(check.GetSummary()); summary != "" {
		output = append(output, "-", summary)
	}

	return o.filterString(warnings, strings.Join(output, " "))
}

func (o *Runtime) buildCheckPerfDataOutput(perfData []*PerfData, separator string, warnings *WarningCollection) (string, error) {
	perfDataStrings := make([]string, len(perfData))
	for key, value := range perfData {
		if output, err := value.BuildOutput(); err == nil {
			perfDataStrings[key] = output
		} else {
			return "", err
		}
	}

	return o.filterString(warnings, strings.Join(perfDataStrings, separator)), nil
}

func (o *Runtime) filterString(warnings *WarningCollection, value string) string {
	originalValue := value
	for _, character := range illegalCharacters {
		value = strings.Replace(value, character, "", -1)
		if originalValue != value && warnings != nil {
			warnings.Add(NewWarning(fmt.Sprintf("nagopher: stripped illegal character from string [%s]",
				originalValue)))
		}
	}

	return value
}

func (o *Runtime) filterStrings(warnings *WarningCollection, values []string) []string {
	results := make([]string, len(values))
	for key, value := range values {
		results[key] = o.filterString(warnings, value)
	}

	return results
}
