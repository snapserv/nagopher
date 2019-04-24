package nagopher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptionalBounds_Empty(t *testing.T) {
	// when
	optionalBounds := OptionalBounds{}
	actualBounds, err := optionalBounds.Get()

	// then
	assert.Error(t, err)
	assert.Empty(t, actualBounds)
}

func TestNewOptionalBounds(t *testing.T) {
	// given
	expectedBounds := NewBounds()

	// when
	optionalBounds := NewOptionalBounds(expectedBounds)
	actualBounds, err := optionalBounds.Get()

	// then
	assert.NoError(t, err)
	assert.Equal(t, true, optionalBounds.Present())
	assert.Equal(t, expectedBounds, actualBounds)
}

func TestOptionalBounds_OrElse(t *testing.T) {
	// given
	expectedBounds := NewBounds()
	alternativeBounds := NewBounds()

	// when
	optionalBounds1 := NewOptionalBounds(expectedBounds)
	optionalBounds2 := OptionalBounds{}

	// then
	assert.Equal(t, expectedBounds, optionalBounds1.OrElse(alternativeBounds))
	assert.Equal(t, alternativeBounds, optionalBounds2.OrElse(alternativeBounds))
}

func TestOptionalBounds_Set(t *testing.T) {
	// given
	bounds1 := NewBounds()
	bounds2 := NewBounds()
	alternativeBounds := NewBounds()
	optionalBounds := NewOptionalBounds(bounds1)

	// when
	optionalBounds.Set(bounds2)

	// then
	assert.Equal(t, true, optionalBounds.Present())
	assert.Equal(t, bounds2, optionalBounds.OrElse(alternativeBounds))
}

func TestOptionalBounds_If(t *testing.T) {
	// given
	var actualBoundsPtr *Bounds = nil
	expectedBounds := NewBounds()
	optionalBounds := NewOptionalBounds(expectedBounds)

	// when
	optionalBounds.If(func(bounds Bounds) {
		actualBoundsPtr = &bounds
	})

	// then
	assert.Equal(t, &expectedBounds, actualBoundsPtr)
}
