package ktaxes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaxCalculate(t *testing.T) {
	// Arrange
	cases := []struct {
		request  TaxRequest
		expected float64
	}{
		{TaxRequest{60000.0, 0.0, nil}, 0.0},
		{TaxRequest{210000.0, 0.0, nil}, 0.0},
		{TaxRequest{500000.0, 0.0, nil}, 29000.0},
		{TaxRequest{560000.0, 0.0, nil}, 35000.0},
		{TaxRequest{1060000.0, 0.0, nil}, 110000.0},
		{TaxRequest{2060000.0, 0.0, nil}, 310000.0},
		{TaxRequest{500000.0, 30000.0, nil}, -1000.0},
		{TaxRequest{560000.0, 36000, nil}, -1000.0},
		{TaxRequest{1060000.0, 111000, nil}, -1000.0},
		{TaxRequest{2060000.0, 311000, nil}, -1000.0},
	}

	// Act & Assert
	for _, c := range cases {
		// Act
		tax := taxCalculate(c.request.TotalIncome, c.request.Wht)
		// Assert
		assert.Equal(t, c.expected, tax, "tax calculation is incorrect : totalIncome = %.2f , wht = %.2f", c.request.TotalIncome, c.request.Wht)
	}
}
