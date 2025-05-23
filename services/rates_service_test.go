package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildRatesMap(t *testing.T) {
	tests := []struct {
		name     string
		rates    map[string]float64
		expected RatesMap
	}{
		{
			name: "2 currencies",
			rates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.5,
			},
			expected: RatesMap{
				"USD": {
					"EUR": 0.5,
				},
				"EUR": {
					"USD": 2,
				},
			},
		},
		{
			name: "3 currencies",
			rates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.5,
				"GBP": 0.8,
			},
			expected: RatesMap{
				"USD": {
					"EUR": 0.5,
					"GBP": 0.8,
				},
				"EUR": {
					"USD": 2,
					"GBP": 1.6,
				},
				"GBP": {
					"USD": 1.25,
					"EUR": 0.625,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rates := buildRatesMap(test.rates)
			assert.Equal(t, test.expected, rates)
		})
	}
}
