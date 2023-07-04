package cart_domain

import "github.com/shopspring/decimal"

const taxRate = 0.10

func calculateSubtotal(products []CartProduct) decimal.Decimal {
	subtotal := decimal.NewFromFloat(0.00)
	for _, product := range products {
		subtotal = subtotal.Add(product.Subtotal)
	}
	return subtotal
}

func calculateShipping() decimal.Decimal {
	return decimal.NewFromFloat(5.00)
}

func calculateTax(subtotal decimal.Decimal) decimal.Decimal {
	return subtotal.Mul(decimal.NewFromFloat(taxRate))
}

func calculateTotal(subtotal decimal.Decimal, shipping decimal.Decimal, tax decimal.Decimal) decimal.Decimal {
	return subtotal.Add(shipping).Add(tax)
}
