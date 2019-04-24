package nagopher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptionalContext_Empty(t *testing.T) {
	// when
	optionalContext := OptionalContext{}
	actualContext, err := optionalContext.Get()

	// then
	assert.Error(t, err)
	assert.Empty(t, actualContext)
}

func TestNewOptionalContext(t *testing.T) {
	// given
	expectedContext := NewStringInfoContext("expected")

	// when
	optionalContext := NewOptionalContext(expectedContext)
	actualContext, err := optionalContext.Get()

	// then
	assert.NoError(t, err)
	assert.Equal(t, true, optionalContext.Present())
	assert.Equal(t, expectedContext, actualContext)
}

func TestOptionalContext_OrElse(t *testing.T) {
	// given
	expectedContext := NewStringInfoContext("expected")
	alternativeContext := NewStringInfoContext("alternative")

	// when
	optionalContext1 := NewOptionalContext(expectedContext)
	optionalContext2 := OptionalContext{}

	// then
	assert.Equal(t, expectedContext, optionalContext1.OrElse(alternativeContext))
	assert.Equal(t, alternativeContext, optionalContext2.OrElse(alternativeContext))
}

func TestOptionalContext_Set(t *testing.T) {
	// given
	context1 := NewStringInfoContext("context 1")
	context2 := NewStringInfoContext("context 2")
	alternativeContext := NewStringInfoContext("alternative")

	// when
	optionalContext := NewOptionalContext(context1)
	optionalContext.Set(context2)

	// then
	assert.Equal(t, true, optionalContext.Present())
	assert.Equal(t, context2, optionalContext.OrElse(alternativeContext))
}

func TestOptionalContext_If(t *testing.T) {
	// given
	var actualContextPtr *Context = nil
	expectedContext := NewStringInfoContext("expected")

	// when
	optionalContext := NewOptionalContext(expectedContext)
	optionalContext.If(func(context Context) {
		actualContextPtr = &context
	})

	// then
	assert.Equal(t, &expectedContext, actualContextPtr)
}
