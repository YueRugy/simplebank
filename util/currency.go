package util

var (
	currencyMap = make(map[string]struct{})
)

func init() {
	for _, t := range currencyType {
		currencyMap[t] = struct{}{}
	}
}

func IsSupportCurrency(currency string) bool {
	_, ok := currencyMap[currency]
	return ok
}
