package util

const (
	USD = "USD"
	NGN = "NGN"
	EUR = "EUR"
)

func IsSupportedCurrency(currencyCode string) bool {
	switch currencyCode {
	case USD, EUR, NGN:
		return true
	}
	return false
}
