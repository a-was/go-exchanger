package services

type RatesMap map[string]map[string]float64

type RatesGetter interface {
	GetRates(targetCurrencies []string) (RatesMap, error)
}

func buildRatesMap(rates map[string]float64) RatesMap {
	ratesMap := make(RatesMap, len(rates))

	for c1, r1 := range rates {
		ratesMap[c1] = make(map[string]float64, len(rates)-1)
		for c2, r2 := range rates {
			if c1 == c2 {
				continue
			}
			ratesMap[c1][c2] = r2 / r1
		}
	}

	return ratesMap
}
