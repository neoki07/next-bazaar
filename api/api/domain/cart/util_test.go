package cart_domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCalculateSubtotal(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		products []CartProduct
		expected decimal.Decimal
	}{
		{
			name:     "empty",
			products: []CartProduct{},
			expected: decimal.NewFromFloat(0.00),
		},
		{
			name: "multiple products",
			products: []CartProduct{
				{Subtotal: decimal.NewFromFloat(10.00)},
				{Subtotal: decimal.NewFromFloat(20.00)},
				{Subtotal: decimal.NewFromFloat(30.00)},
			},
			expected: decimal.NewFromFloat(60.00),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			subtotal := calculateSubtotal(tc.products)

			require.True(t, subtotal.Equal(tc.expected))
		})
	}
}

func TestCalculateShipping(t *testing.T) {
	t.Parallel()

	shipping := calculateShipping()

	require.True(t, shipping.Equal(decimal.NewFromFloat(5.00)))
}

func TestCalculateTax(t *testing.T) {
	t.Parallel()

	subtotal := decimal.NewFromFloat(50.00)

	tax := calculateTax(subtotal)

	require.True(t, tax.Equal(decimal.NewFromFloat(5.00)))
}

func TestCalculateTotal(t *testing.T) {
	t.Parallel()

	subtotal := decimal.NewFromFloat(50.00)
	shipping := decimal.NewFromFloat(5.00)
	tax := decimal.NewFromFloat(5.00)

	total := calculateTotal(subtotal, shipping, tax)

	require.True(t, total.Equal(decimal.NewFromFloat(60.00)))
}
