package util

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
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

func RandomMoney() (string, error) {
	amountWithCurrency, err := faker.GetPrice().AmountWithCurrency(reflect.Value{})
	if err != nil {
		return "", err
	}

	amountWithCurrencyStr, ok := amountWithCurrency.(string)
	if !ok {
		return "", fmt.Errorf("cannot convert price to string")
	}

	return strings.Split(amountWithCurrencyStr, " ")[1], nil
}

func RandomPrice() (string, error) {
	return RandomMoney()
}

func RandomImageUrl() string {
	return "https://picsum.photos/300/300/?random"
}
