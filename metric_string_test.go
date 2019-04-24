package nagopher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStringMetric(t *testing.T) {
	// when
	metric1, err1 := NewStringMetric("string", "", "")
	metric2, err2 := NewStringMetric("", "", "")

	// then
	assert.NoError(t, err1)
	assert.Error(t, err2)
	assert.Implements(t, (*StringMetric)(nil), metric1)
	assert.Nil(t, metric2)
}

func TestMustNewStringMetric(t *testing.T) {
	assert.NotPanics(t, func() {
		MustNewStringMetric("valid", "", "")
	})

	assert.Panics(t, func() {
		MustNewStringMetric("", "", "")
	})
}

func TestStringMetric_ToNagiosValue(t *testing.T) {
	// given
	var value string = "Hello World"

	// when
	metric := MustNewStringMetric("test", value, "")

	// then
	assert.Equal(t, value, metric.ToNagiosValue())
}

func TestStringMetric_Value(t *testing.T) {
	// given
	var value string = "Hello World"

	// when
	metric := MustNewStringMetric("test", value, "")

	// then
	assert.Equal(t, value, metric.Value())
}
