package nagopher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScalarContext_Evaluate(t *testing.T) {
	// given
	warningThreshold := NewBounds(LowerBound(10), UpperBound(20))
	criticalThreshold := NewBounds(LowerBound(0), UpperBound(30))
	context := NewScalarContext("context", &warningThreshold, &criticalThreshold)
	metric1 := MustNewNumericMetric("metric 1", 15, "", nil, "")
	metric2 := MustNewNumericMetric("metric 2", 5, "", nil, "")
	metric3 := MustNewNumericMetric("metric 3", -5, "", nil, "")
	metric4 := MustNewStringMetric("invalid", "Oops!", "")
	resource := NewResource()

	// when
	result1 := context.Evaluate(metric1, resource)
	result2 := context.Evaluate(metric2, resource)
	result3 := context.Evaluate(metric3, resource)
	result4 := context.Evaluate(metric4, resource)

	// then
	assert.Equal(t, StateOk(), result1.State().OrElse(nil))
	assert.Equal(t, StateWarning(), result2.State().OrElse(nil))
	assert.Equal(t, StateCritical(), result3.State().OrElse(nil))
	assert.Equal(t, StateUnknown(), result4.State().OrElse(nil))

	assert.Equal(t, "", result1.Hint())
	assert.Equal(t, "outside range 10:20", result2.Hint())
	assert.Equal(t, "outside range 0:30", result3.Hint())
	assert.Contains(t, result4.Hint(), "ScalarContext can not process metric of type")
}

func TestScalarContext_Performance(t *testing.T) {
	// given
	context := NewScalarContext("context", nil, nil)
	metric1 := MustNewNumericMetric("valid", 42, "", nil, "")
	metric2 := MustNewNumericMetric("inv='alid", 42, "", nil, "")
	resource := NewResource()

	// when
	var perfData1, perfData2 PerfData = nil, nil
	optionalPerfData1, err1 := context.Performance(metric1, resource)
	optionalPerfData2, err2 := context.Performance(metric2, resource)

	if err1 == nil {
		perfData1, err1 = optionalPerfData1.Get()
	}
	if err2 == nil {
		perfData2, err2 = optionalPerfData2.Get()
	}

	// then
	assert.NoError(t, err1)
	assert.Error(t, err2)

	assert.Implements(t, (*PerfData)(nil), perfData1)
	assert.Nil(t, perfData2)
}
