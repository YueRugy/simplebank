package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet   = "abcdefghijklmnopqrstuvwxyz"
	ownerLen   = 6
	balanceMin = 10
	balanceMax = 1000
)

var (
	currencyType = [3]string{"EUR", "USD", "CAD"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[rand.Intn(k)%len(alphabet)])
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(ownerLen)
}

func RandomBalance() int64 {
	return RandomInt(balanceMin, balanceMax)
}

func RandomCurrency() string {
	return currencyType[rand.Intn(len(currencyType))]
}

func RandomAmount(number int64) int64 {
	return RandomInt(0, number)
}
