package nagopher

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWarning(t *testing.T) {
	// when
	formatString := "Value %s %d"
	formatArgs := []interface{}{"is", 13.37}
	warning := NewWarning(formatString, formatArgs...)

	// then
	assert.Equal(t, fmt.Sprintf(formatString, formatArgs...), warning.Warning())
}

func TestWarningCollection_Add(t *testing.T) {
	// given
	warning1 := NewWarning("Hello")
	warning2 := NewWarning("World")
	warnings := NewWarningCollection()

	// when
	warnings.Add(warning1, warning2)

	// then
	assert.Equal(t, 2, len(warnings.Get()))
	assert.Contains(t, warnings.Get(), warning1)
	assert.Contains(t, warnings.Get(), warning2)
}

func TestWarningCollection_GetWarningStrings(t *testing.T) {
	// given
	warning1 := NewWarning("Hello")
	warning2 := NewWarning("World")
	warnings := NewWarningCollection()
	warnings.Add(warning1, warning2)

	// when
	warningStrings := warnings.GetWarningStrings()

	// then
	assert.Equal(t, 2, len(warningStrings))
	assert.Contains(t, warningStrings, warning1.Warning())
	assert.Contains(t, warningStrings, warning2.Warning())
}
