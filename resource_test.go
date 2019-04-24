package nagopher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseResource_Probe(t *testing.T) {
	// given
	resource := NewResource()
	warnings := NewWarningCollection()

	// when
	metrics, err := resource.Probe(warnings)

	// then
	assert.NoError(t, err)
	assert.Empty(t, warnings.Get())
	assert.Empty(t, metrics)
}
