package nagopher

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type MockResource struct {
	Resource
}

func TestBaseRuntime_Execute(t *testing.T) {
	// given
	warningThreshold := NewBounds(LowerBound(10), UpperBound(80))
	check1 := NewCheck("usage", NewSummarizer())
	check2 := NewCheck("usage", NewSummarizer())

	check1.AttachResources(NewMockResource())
	check1.AttachContexts(NewScalarContext("usage", nil, nil))
	check2.AttachResources(NewMockResource())
	check2.AttachContexts(NewScalarContext("usage", &warningThreshold, nil))

	// when
	actualResult1 := NewRuntime(false).Execute(check1) // non-verbose
	actualResult2 := NewRuntime(true).Execute(check2)  // verbose

	// then
	assert.Equal(t, StateOk().ExitCode(), actualResult1.ExitCode())
	assert.Equal(t, strings.Join([]string{
		"USAGE OK - usage1 is 49.4% | usage1=49.4% usage2=92.6% usage3=83.1",
		"nagopher: stripped illegal character from string [usage3=83.1]",
	}, "\n")+"\n", actualResult1.Output())

	assert.Equal(t, StateWarning().ExitCode(), actualResult2.ExitCode())
	assert.Equal(t, strings.Join([]string{
		"USAGE WARNING - usage2 is 92.6% (outside range 10:80) | usage1=49.4%;10:80 usage2=92.6%;10:80 usage3=83.1;10:80",
		"warning: usage2 is 92.6% (outside range 10:80)",
		"warning: usage3 is 83.1 (outside range 10:80)",
		"nagopher: stripped illegal character from string [usage3=83.1;10:80]",
		"nagopher: stripped illegal character from string [warning: usage3 is 83.1 (outside range 10:80)]",
	}, "\n")+"\n", actualResult2.Output())
}

func NewMockResource() Resource {
	return &MockResource{
		Resource: NewResource(),
	}
}

func (r MockResource) Probe(warnings WarningCollection) ([]Metric, error) {
	return []Metric{
		MustNewNumericMetric("usage1", 49.4, "%", nil, "usage"),
		MustNewNumericMetric("usage2", 92.6, "%", nil, "usage"),
		MustNewNumericMetric("usage3", 83.1, "|", nil, "usage"),
	}, nil
}

func TestNewCheckResult(t *testing.T) {
	// when
	exitCode := StateOk().ExitCode()
	description := StateOk().Description()
	checkResult := NewCheckResult(exitCode, description)

	// then
	assert.Equal(t, exitCode, checkResult.ExitCode())
	assert.Equal(t, description, checkResult.Output())
}
