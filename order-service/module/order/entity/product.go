package entity

import "github.com/shopspring/decimal"

type Product struct {
	ID    string
	Name  string
	Price decimal.Decimal
}
