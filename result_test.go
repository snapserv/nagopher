package nagopher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewResult(t *testing.T) {
	// given
	expectedState := StateOk()
	expectedMetric := MustNewStringMetric("metric", "test", "context")
	expectedContext := NewStringInfoContext("context")
	expectedResource := NewResource()
	expectedHint := "Result Hint"

	// when
	result := NewResult(
		ResultState(expectedState), ResultHint(expectedHint),
		ResultMetric(expectedMetric), ResultContext(expectedContext), ResultResource(expectedResource),
	)
	actualState, _ := result.State().Get()
	actualMetric, _ := result.Metric().Get()
	actualContext, _ := result.Context().Get()
	actualResource, _ := result.Resource().Get()

	// then
	assert.Equal(t, expectedHint, result.Hint())
	assert.Equal(t, expectedState, actualState)
	assert.Equal(t, expectedMetric, actualMetric)
	assert.Equal(t, expectedContext, actualContext)
	assert.Equal(t, expectedResource, actualResource)
}

func TestResult_String(t *testing.T) {
	// given
	metric := MustNewNumericMetric("metric", 13.37, "", nil, "")
	context := NewScalarContext("context", nil, nil)
	hint := "Result Hint"

	// when
	result1 := NewResult(ResultContext(context), ResultMetric(metric))
	result2 := NewResult(ResultMetric(metric))
	result3 := NewResult(ResultContext(context), ResultMetric(metric), ResultHint(hint))
	result4 := NewResult(ResultMetric(metric), ResultHint(hint))

	// then
	assert.Equal(t, "metric is 13.37", result1.String())
	assert.Equal(t, "13.37", result2.String())
	assert.Equal(t, "metric is 13.37 (Result Hint)", result3.String())
	assert.Equal(t, "13.37 (Result Hint)", result4.String())
}
