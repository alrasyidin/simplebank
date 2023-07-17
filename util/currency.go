package util

var Currencies = []string{"EUR", "USD", "IDR"}

func IsValidCurrency(currency string) bool {

	for _, v := range Currencies {
		if v == currency {
			return true
		}
	}

	return false
}
