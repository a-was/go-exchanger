package services

type MockRatesService struct {
}

var _ RatesGetter = (*MockRatesService)(nil)

func (m *MockRatesService) GetRates(targetCurrencies []string) (map[string]float64, error) {
	return nil, nil
}
