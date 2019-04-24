package nagopher

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNewNumericMetric(t *testing.T) {
	// when
	metric1, err1 := NewNumericMetric("valid", 13.37, "K", nil, "")
	metric2, err2 := NewNumericMetric("", 0, "", nil, "")

	// then
	assert.NoError(t, err1)
	assert.Error(t, err2)
	assert.Implements(t, (*NumericMetric)(nil), metric1)
	assert.Nil(t, metric2)
}

func TestMustNewNumericMetric(t *testing.T) {
	assert.NotPanics(t, func() {
		MustNewNumericMetric("valid", 13.37, "K", nil, "")
	})

	assert.Panics(t, func() {
		MustNewNumericMetric("", 0, "", nil, "")
	})
}

func TestNumericMetric_ToNagiosValue(t *testing.T) {
	// given
	metric1 := MustNewNumericMetric("metric 1", 10, "B", nil, "")
	metric2 := MustNewNumericMetric("metric 2", 13.37, "B", nil, "")
	metric3 := MustNewNumericMetric("metric 3", math.NaN(), "", nil, "")

	// when
	value1 := metric1.ToNagiosValue()
	value2 := metric2.ToNagiosValue()
	value3 := metric3.ToNagiosValue()

	// then
	assert.Equal(t, "10B", value1)
	assert.Equal(t, "13.37B", value2)
	assert.Equal(t, "U", value3)
}

func TestNumericMetric_Value(t *testing.T) {
	// given
	var value float64 = 13.37

	// when
	metric := MustNewNumericMetric("metric", value, "", nil, "")

	// then
	assert.Equal(t, value, metric.Value())
}
