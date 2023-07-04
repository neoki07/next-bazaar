package db

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestDecimalMarshalJSON(t *testing.T) {
	d := Decimal{decimal.NewFromFloat(123.456)}

	b, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, `"123.456"`, string(b))
}

func TestDecimalUnmarshalJSON(t *testing.T) {
	var d Decimal

	err := json.Unmarshal([]byte(`"123.456"`), &d)
	if err != nil {
		t.Fatal(err)
	}

	require.True(t, decimal.NewFromFloat(123.456).Equal(d.Decimal))
}
