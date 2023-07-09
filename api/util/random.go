package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt32 generates a random int32
func RandomInt32(n int32) int32 {
	return rand.Int31n(n)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomName generates a random name
func RandomName() string {
	return RandomString(6)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// RandomUUID generates a uuid
func RandomUUID() uuid.UUID {
	return uuid.New()
}

// RandomPrice generates a random price
func RandomPrice() decimal.Decimal {
	min := 1.00
	max := 100.00
	randomFloat := min + rand.Float64()*(max-min)
	return decimal.NewFromFloat(randomFloat).Round(2)
}
