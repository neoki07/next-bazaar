package db

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type Decimal struct {
	decimal.Decimal
}

func (d Decimal) MarshalJSON() ([]byte, error) {

	return json.Marshal(d.String())
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	dec, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}

	d.Decimal = dec
	return nil
}
