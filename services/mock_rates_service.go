package services

type MockRatesService struct {
	RatesMap RatesMap
}

var _ RatesGetter = (*MockRatesService)(nil)

func (m *MockRatesService) GetRates(targetCurrencies []string) (RatesMap, error) {
	return m.RatesMap, nil
}
