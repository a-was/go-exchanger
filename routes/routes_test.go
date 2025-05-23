package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-was/go-exchanger/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_GetRates(t *testing.T) {

	tests := []struct {
		name         string
		query        string
		rates        map[string]float64
		expectedCode int
		expectedBody []outRate
	}{
		{
			name:  "2 currencies",
			query: "USD,EUR",
			rates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.5,
			},
			expectedCode: http.StatusOK,
			expectedBody: []outRate{
				{From: "USD", To: "EUR", Rate: 0.5},
				{From: "EUR", To: "USD", Rate: 2},
			},
		},
		{
			name:  "3 currencies",
			query: "USD,EUR,PLN",
			rates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.5,
				"PLN": 4,
			},
			expectedCode: http.StatusOK,
			expectedBody: []outRate{
				{From: "USD", To: "EUR", Rate: 0.5},
				{From: "USD", To: "PLN", Rate: 4},
				{From: "EUR", To: "USD", Rate: 2},
				{From: "EUR", To: "PLN", Rate: 8},
				{From: "PLN", To: "USD", Rate: 0.25},
				{From: "PLN", To: "EUR", Rate: 0.125},
			},
		},
		{
			name:         "empty query",
			query:        "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid query",
			query:        "abc",
			expectedCode: http.StatusBadRequest,
		},
	}

	r := Router{
		Engine: gin.Default(),
	}
	r.RegisterRoutes()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r.RatesService = &services.MockRatesService{
				RatesMap: services.BuildRatesMap(test.rates),
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/rates?currencies="+test.query, nil)
			r.Engine.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)

			if len(test.expectedBody) > 0 {
				var response []outRate
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, len(test.expectedBody), len(response))
				for _, rate := range response {
					assert.Contains(t, test.expectedBody, rate)
				}
			}
		})
	}
}
