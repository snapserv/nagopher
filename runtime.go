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

// Runtime executes a specific Check instance and prints or outputs the results according to the Nagios plugin specs
type Runtime interface {
	Execute(Check) CheckResult
	ExecuteAndExit(Check)
}

// CheckResult contains the results of a Check together with an exit code to indicate the check state
type CheckResult interface {
	ExitCode() int8
	Output() string
}

type baseRuntime struct {
	verboseOutput bool
}

type checkResult struct {
	exitCode int8
	output   string
}

var illegalOutputChars = []string{"|"}

// NewRuntime instantiates a new Runtime, optionally enabling verbose output
func NewRuntime(verboseOutput bool) Runtime {
	runtime := &baseRuntime{
		verboseOutput: verboseOutput,
	}

	return runtime
}

func (r baseRuntime) Execute(check Check) CheckResult {
	warnings := NewWarningCollection()
	check.Run(warnings)

	checkState := check.State()
	checkOutput := r.buildNagiosOutput(check, warnings)

	return NewCheckResult(checkState.ExitCode(), checkOutput)
}

func (r baseRuntime) ExecuteAndExit(check Check) {
	result := r.Execute(check)
	fmt.Print(result.Output())
	os.Exit(int(result.ExitCode()))
}

func (r baseRuntime) buildNagiosOutput(check Check, warnings WarningCollection) string {
	var outputParts []string

	outputParts = append(outputParts, r.buildNagiosStatus(check, warnings))
	if perfData := r.buildNagiosPerfData(check.PerfData(), warnings); perfData != "" {
		outputParts = append(outputParts, " | ", perfData)
	}
	outputParts = append(outputParts, "\n")

	if r.verboseOutput {
		lines := r.sanitizeStrings(check.VerboseSummary(), warnings)
		if len(lines) > 0 {
			outputParts = append(outputParts, strings.Join(lines, "\n"), "\n")
		}
	}

	if warningStrings := warnings.GetWarningStrings(); len(warningStrings) > 0 {
		warningStrings = r.sanitizeStrings(warningStrings, nil)
		outputParts = append(outputParts, strings.Join(warningStrings, "\n"), "\n")
	}

	return strings.Join(outputParts, "")
}

func (r baseRuntime) buildNagiosStatus(check Check, warnings WarningCollection) string {
	var outputParts []string

	if check.Name() != "" {
		outputParts = append(outputParts, strings.ToUpper(check.Name()))
	}

	outputParts = append(outputParts, strings.ToUpper(check.State().Description()))
	summary := strings.TrimSpace(check.Summary())
	if summary != "" {
		outputParts = append(outputParts, "-", summary)
	}

	return strings.Join(r.sanitizeStrings(outputParts, warnings), " ")
}

func (r baseRuntime) buildNagiosPerfData(perfData []PerfData, warnings WarningCollection) string {
	outputParts := make([]string, len(perfData))
	for key, value := range perfData {
		outputParts[key] = value.ToNagiosPerfData()
	}

	return strings.Join(r.sanitizeStrings(outputParts, warnings), " ")
}

func (r baseRuntime) sanitizeStrings(values []string, warnings WarningCollection) []string {
	results := make([]string, len(values))
	for key, value := range values {
		results[key] = r.sanitizeString(value, warnings)
	}

	return results
}

func (r baseRuntime) sanitizeString(value string, warnings WarningCollection) string {
	originalValue := value
	for _, character := range illegalOutputChars {
		value = strings.Replace(value, character, "", -1)
		if originalValue != value && warnings != nil {
			warnings.Add(NewWarning("nagopher: stripped illegal character from string [%s]", originalValue))
		}
	}

	return value
}

// NewCheckResult instantiates a new CheckResult with the given exit code and output string
func NewCheckResult(exitCode int8, output string) CheckResult {
	checkResult := &checkResult{
		exitCode: exitCode,
		output:   output,
	}

	return checkResult
}

func (r checkResult) ExitCode() int8 {
	return r.exitCode
}

func (r checkResult) Output() string {
	return r.output
}
