package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-was/go-exchanger/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_GetRates(t *testing.T) {

	r := Router{
		Engine:       gin.Default(),
		RatesService: &services.MockRatesService{},
	}
	r.RegisterRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rates?currencies=GBP,EUR", nil)
	r.Engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_GetCurrenciesMap(t *testing.T) {
	tests := []struct {
		name       string
		currencies []string
		expected   map[string][]string
	}{
		{
			name:       "2 currencies",
			currencies: []string{"USD", "EUR"},
			expected: map[string][]string{
				"USD": {"EUR"},
				"EUR": {"USD"},
			},
		},
		{
			name:       "3 currencies",
			currencies: []string{"USD", "EUR", "GBP"},
			expected: map[string][]string{
				"USD": {"EUR", "GBP"},
				"EUR": {"USD", "GBP"},
				"GBP": {"USD", "EUR"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rates := GetCurrenciesMap(test.currencies)
			assert.Equal(t, test.expected, rates)
		})
	}
}
