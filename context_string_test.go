package nagopher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringInfoContext_Evaluate(t *testing.T) {
	// given
	context := NewStringInfoContext("context")
	stringInfoContext, _ := context.(*stringInfoContext)
	metric := MustNewStringMetric("test", "Hello World", "")
	resource := NewResource()

	// when
	result := context.Evaluate(metric, resource)
	resultState, _ := result.State().Get()
	resultMetric, _ := result.Metric().Get()
	resultResource, _ := result.Resource().Get()
	resultContext, _ := result.Context().Get()

	// then
	assert.Equal(t, StateOk(), resultState)
	assert.Equal(t, metric, resultMetric)
	assert.Equal(t, resource, resultResource)
	assert.Equal(t, *stringInfoContext, resultContext)
}

func TestStringMatchContext_Evaluate(t *testing.T) {
	// given
	context := NewStringMatchContext("context", StateWarning(), []string{"ABC", "def", "gHi"})
	metric1 := MustNewStringMetric("metric 1", "ABC", "")
	metric2 := MustNewStringMetric("metric 2", "dEf", "")
	metric3 := MustNewStringMetric("metric 3", "GHj", "")
	metric4 := MustNewNumericMetric("invalid", 42, "", nil, "")
	resource := NewResource()

	// when
	result1 := context.Evaluate(metric1, resource)
	result2 := context.Evaluate(metric2, resource)
	result3 := context.Evaluate(metric3, resource)
	result4 := context.Evaluate(metric4, resource)

	// then
	assert.Equal(t, StateOk(), result1.State().OrElse(nil))
	assert.Equal(t, StateOk(), result2.State().OrElse(nil))
	assert.Equal(t, StateWarning(), result3.State().OrElse(nil))
	assert.Equal(t, StateUnknown(), result4.State().OrElse(nil))

	assert.Equal(t, "", result1.Hint())
	assert.Equal(t, "", result2.Hint())
	assert.Equal(t, "got [ghj], expected [abc],[def],[ghi]", result3.Hint())
	assert.Contains(t, result4.Hint(), "StringMatchContext can not process metric of type")
}