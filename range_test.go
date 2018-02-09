package nagopher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRange_String_Empty(t *testing.T) {
	assert.Equal(t, "", ParseRange("").String())
}

func TestRange_String_ExplicitStartEnd(t *testing.T) {
	assert.Equal(t, "1:9", ParseRange("1:9").String())
}

func TestRange_String_OmitStart(t *testing.T) {
	assert.Equal(t, "9", ParseRange("9").String())
}

func TestRange_String_OmitEnd(t *testing.T) {
	assert.Equal(t, "1:", ParseRange("1:").String())
}

func TestRange_String_NegativeInfinityStart(t *testing.T) {
	assert.Equal(t, "~:10", ParseRange("~:10").String())
}

func TestRange_String_NegativeInfinityEnd(t *testing.T) {
	// TODO: Implement test case
	// assert.Equal(t, "10:~", ParseRange("10:~").String())
}

func TestRange_String_Invert(t *testing.T) {
	assert.Equal(t, "@1:9", ParseRange("@1:9").String())
}

func TestRange_String_LargeNumberStart(t *testing.T) {
	assert.Equal(t, "4200000000:", ParseRange("4200000000:").String())
}

func TestRange_String_LargeNumberEnd(t *testing.T) {
	assert.Equal(t, "4200000000", ParseRange("4200000000").String())
}

func TestRange_ViolationHint_Normal(t *testing.T) {
	assert.Equal(t, "outside range 1:9", ParseRange("1:9").ViolationHint())
}

func TestRange_ViolationHint_OmitStart(t *testing.T) {
	assert.Equal(t, "outside range 0:9", ParseRange(":9").ViolationHint())
}

func TestRange_ViolationHint_OmitEnd(t *testing.T) {
	assert.Equal(t, "outside range 1:inf", ParseRange("1:").ViolationHint())
}

func TestRange_ViolationHint_NegativeInfinityStart(t *testing.T) {
	assert.Equal(t, "outside range -inf:1", ParseRange("~:1").ViolationHint())
}

func TestRange_ViolationHint_Empty(t *testing.T) {
	assert.Equal(t, "outside range 0:inf", ParseRange("").ViolationHint())
}
