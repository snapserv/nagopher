package nagopher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseSummarizer_Empty(t *testing.T) {
	// when
	summarizer := NewSummarizer()

	// then
	assert.Equal(t, "No check results", summarizer.Empty())
}

func TestBaseSummarizer_Ok(t *testing.T) {
	// given
	summarizer := NewSummarizer()
	check := NewCheck("check", summarizer)

	// when
	check.Results().Add(
		NewResult(ResultState(StateOk()), ResultHint("Result 1")),
		NewResult(ResultState(StateOk()), ResultHint("Result 2")),
	)

	// then
	assert.Equal(t, "Result 1", summarizer.Ok(check))
}

func TestBaseSummarizer_Fallback(t *testing.T) {
	// given
	summarizer := NewSummarizer()

	// when
	check := NewCheck("check", summarizer)

	// then
	assert.Equal(t, summarizer.Empty(), summarizer.Ok(check))
	assert.Equal(t, summarizer.Empty(), summarizer.Problem(check))
}

func TestBaseSummarizer_Problem(t *testing.T) {
	// given
	summarizer := NewSummarizer()
	check := NewCheck("check", summarizer)

	// when
	check.Results().Add(
		NewResult(ResultState(StateOk()), ResultHint("Result OK")),
		NewResult(ResultState(StateCritical()), ResultHint("Result CRITICAL")),
		NewResult(ResultState(StateWarning()), ResultHint("Result WARNING")),
		NewResult(ResultState(StateOk()), ResultHint("Result OK")),
	)

	// then
	assert.Equal(t, "Result CRITICAL", summarizer.Problem(check))
}

func TestBaseSummarizer_Verbose(t *testing.T) {
	// given
	summarizer := NewSummarizer()
	check := NewCheck("check", summarizer)

	// when
	check.Results().Add(
		NewResult(ResultState(StateOk()), ResultHint("hidden from output")),
		NewResult(ResultState(StateWarning()), ResultHint("Reason 1")),
		NewResult(ResultState(StateCritical()), ResultHint("Reason 2")),
		NewResult(ResultState(StateCritical()), ResultHint("Reason 3")),
		NewResult(ResultHint("Informational Result")),
	)

	// then
	expected := []string{"critical: Reason 2", "critical: Reason 3", "warning: Reason 1", "info: Informational Result"}
	assert.Equal(t, expected, summarizer.Verbose(check))
}
