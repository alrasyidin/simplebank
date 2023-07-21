package util

const (
	USD = "USD"
	IDR = "IDR"
	EUR = "EUR"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case USD, EUR, IDR:
		return true
	}

	return false
}
