package nagopher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPerfData_String_NormalLabel(t *testing.T) {
	perfData := NewPerfData("i", 42, "", nil, nil, nil)
	assert.Equal(t, "i=42", perfData.String())
}

func TestPerfData_String_QuotedLabel(t *testing.T) {
	perfData := NewPerfData("i i", 42, "", nil, nil, nil)
	assert.Equal(t, "'i i'=42", perfData.String())
}

func TestPerfData_String_InvalidLabel_Quotes(t *testing.T) {
	// TODO: Implement test case
}

func TestPerfData_String_InvalidLabel_Equals(t *testing.T) {
	// TODO: Implement test case
}

func TestPerfData_String_ValueUnit(t *testing.T) {
	perfData := NewPerfData("i", 42, "%", nil, nil, nil)
	assert.Equal(t, "i=42%", perfData.String())
}

func TestPerfData_String_ValueRange(t *testing.T) {
	valueRange := NewRange(false, 0, 100)
	perfData := NewPerfData("i", 42, "", valueRange, nil, nil)

	assert.Equal(t, "i=42;;;0;100", perfData.String())
}

func TestPerfData_String_WarningRange(t *testing.T) {
	warningRange := NewRange(false, 0, 100)
	perfData := NewPerfData("i", 42, "", nil, warningRange, nil)

	assert.Equal(t, "i=42;0:100", perfData.String())
}

func TestPerfData_String_CriticalRange(t *testing.T) {
	criticalRange := NewRange(false, 0, 100)
	perfData := NewPerfData("i", 42, "", nil, nil, criticalRange)

	assert.Equal(t, "i=42;;0:100", perfData.String())
}
