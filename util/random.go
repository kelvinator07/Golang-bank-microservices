package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
	numbers  = "0123456789"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// RandomAccountName generates a random owner name
func RandomAccountName() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "NGN"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomPhoneNumber() int64 {
	v, _ := strconv.Atoi(stringWithCharset(11, numbers))
	return int64(v)
}

func RandomAccountNumber() int64 {
	v, _ := strconv.Atoi(stringWithCharset(10, numbers))
	return int64(v)
}

func RandomEmail() string {
	return fmt.Sprintf("%v@email.com", RandomString(6))
}

// RandomStatus generates a random account status: active or inactive
func RandomStatus() string {
	statuses := []string{"active", "inactive"}
	n := len(statuses)
	return statuses[rand.Intn(n)]
}

// RandomGender generates a random gender
func RandomGender() string {
	genders := []string{"MALE", "FEMALE"}
	n := len(genders)
	return genders[rand.Intn(n)]
}
