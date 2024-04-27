package ktaxes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaxCalculate(t *testing.T) {
	// Arrange
	cases := []struct {
		totalIncome float64
		expected    float64
	}{
		{60000.0, 0.0},
		{210000.0, 0.0},
		{500000.0, 29000.0},
		{560000.0, 35000.0},
		{1060000.0, 110000.0},
		{2060000.0, 310000.0}}

	// Act & Assert
	for _, c := range cases {
		// Act
		tax := taxCalculate(c.totalIncome)
		// Assert
		assert.Equal(t, c.expected, tax, "tax calculation is incorrect : totalIncome = %.2f", c.totalIncome)
	}
}
