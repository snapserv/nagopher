package nagopher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetric_ValueUnit_Float(t *testing.T) {
	metric := NewMetric("time", 3.141, "s", nil, "")
	assert.Equal(t, "3.141s", metric.ValueUnit())
}

func TestMetric_ValueUnit_Integer(t *testing.T) {
	metric := NewMetric("count", 42, "", nil, "")
	assert.Equal(t, "42", metric.ValueUnit())
}

func TestMetric_ValueUnit_LargeInteger(t *testing.T) {
	metric := NewMetric("grains", 4200000000, "", nil, "")
	assert.Equal(t, "4200000000", metric.ValueUnit())
}

func TestMetric_ValueUnit_LargeFloat(t *testing.T) {
	metric := NewMetric("grains", 420000.42, "", nil, "")
	assert.Equal(t, "420000.42", metric.ValueUnit())
}
