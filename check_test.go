package nagopher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseCheck_GetSetMeta(t *testing.T) {
	// given
	check := NewCheck("check", NewSummarizer())
	value1 := float64(13.37)
	value2 := "Hello World"

	// when
	check.SetMeta("test-1", value1)
	check.SetMeta("test-2", value2)

	// then
	assert.Equal(t, value1, check.GetMeta("test-1", nil))
	assert.Equal(t, value2, check.GetMeta("test-2", nil))
	assert.Nil(t, check.GetMeta("missing", nil))
}

func TestBaseCheck_AttachContexts(t *testing.T) {
	// given
	context1 := NewStringInfoContext("context 1")
	context2 := NewStringInfoContext("context 2")
	check := NewCheck("check", NewSummarizer())

	// when
	check.AttachContexts(context1)
	check.AttachContexts(context1, context2)

	// then
	assert.Equal(t, 2, len(check.Contexts()))
	assert.Contains(t, check.Contexts(), context1)
	assert.Contains(t, check.Contexts(), context2)
}

func TestBaseCheck_AttachResources(t *testing.T) {
	// given
	resource1 := NewResource()
	resource2 := NewResource()
	check := NewCheck("check", NewSummarizer())

	// when
	check.AttachResources(resource1)
	check.AttachResources(resource1, resource2)

	// then
	assert.Equal(t, 2, len(check.Resources()))
	assert.Contains(t, check.Resources(), resource1)
	assert.Contains(t, check.Resources(), resource2)
}
