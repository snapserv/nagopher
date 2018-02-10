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
	warnings := newWarningCollection()
	check.Run(warnings)

	err, checkOutput := o.buildCheckOutput(warnings, check)
	if err != nil {
		panic(fmt.Sprintf("nagopher: unexpected runtime error [%s]", err.Error()))
	}

	return CheckResult{
		ExitCode: check.GetState().ExitCode,
		Output:   checkOutput,
	}
}

func (o *Runtime) ExecuteAndExit(check *Check) {
	result := o.Execute(check)
	fmt.Print(result.Output)
	os.Exit(result.ExitCode)
}

func (o *Runtime) buildCheckOutput(warnings *warningCollection, check *Check) (error, string) {
	output := o.buildCheckStatusOutput(warnings, check)

	if err, perfData := o.buildCheckPerfDataOutput(warnings, check.GetPerfData(), " "); err == nil {
		output += " | " + perfData
	} else {
		return err, ""
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

	return nil, output
}

func (o *Runtime) buildCheckStatusOutput(warnings *warningCollection, check *Check) string {
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

func (o *Runtime) buildCheckPerfDataOutput(warnings *warningCollection, perfData []*PerfData, separator string) (error, string) {
	perfDataStrings := make([]string, len(perfData))
	for key, value := range perfData {
		if err, output := value.BuildOutput(); err == nil {
			perfDataStrings[key] = output
		} else {
			return err, ""
		}
	}

	return nil, o.filterString(warnings, strings.Join(perfDataStrings, separator))
}

func (o *Runtime) filterString(warnings *warningCollection, value string) string {
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

func (o *Runtime) filterStrings(warnings *warningCollection, values []string) []string {
	results := make([]string, len(values))
	for key, value := range values {
		results[key] = o.filterString(warnings, value)
	}

	return results
}
