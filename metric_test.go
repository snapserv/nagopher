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
