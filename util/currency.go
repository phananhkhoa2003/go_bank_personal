package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	VND = "VND"
	JPY = "JPY"
	AUD = "AUD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, VND, JPY, AUD:
		return true
	default:
		return false
	}
}
