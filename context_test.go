package nagopher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseContext_Describe(t *testing.T) {
	// given
	metric := MustNewNumericMetric("test", 13.37, "apples", nil, "")
	context := newBaseContext("Test Context", "name=%<name>s value=%<value>s unit=%<unit>s")

	// when
	description := context.Describe(metric)

	// then
	assert.Equal(t, "name=test value=13.37 unit=apples", description)
}

func TestBaseContext_Evaluate(t *testing.T) {
	// given
	context := newBaseContext("Test Context", "%<value>s")
	expectedMetric := MustNewStringMetric("string", "Hello", "")
	expectedResource := NewResource()

	// when
	result := context.Evaluate(expectedMetric, expectedResource)
	resultState, _ := result.State().Get()
	resultContext, _ := result.Context().Get()
	resultMetric, _ := result.Metric().Get()
	resultResource, _ := result.Resource().Get()

	// then
	assert.Equal(t, true, result.State().Present())
	assert.Equal(t, true, result.Context().Present())
	assert.Equal(t, true, result.Metric().Present())
	assert.Equal(t, true, result.Resource().Present())

	assert.Equal(t, StateOk(), resultState)
	assert.Equal(t, *context, resultContext)
	assert.Equal(t, expectedMetric, resultMetric)
	assert.Equal(t, expectedResource, resultResource)
}

func TestBaseContext_Performance(t *testing.T) {
	// given
	context := newBaseContext("Test Context", "%<value>s")
	metric := MustNewStringMetric("string", "Hello", "")
	resource := NewResource()

	// when
	perfData, err := context.Performance(metric, resource)

	// then
	assert.NoError(t, err)
	assert.Equal(t, OptionalPerfData{}, perfData)
}
