package services

type RatesGetter interface {
	GetRates(targetCurrencies []string) (map[string]float64, error)
}
