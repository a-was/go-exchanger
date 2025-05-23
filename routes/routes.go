package routes

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/a-was/go-exchanger/services"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine       *gin.Engine
	RatesService services.RatesGetter
}

func (r *Router) RegisterRoutes() {
	r.Engine.GET("/rates", r.getRates)
}

func (r *Router) getRates(c *gin.Context) {
	currenciesQ := c.Query("currencies")
	if currenciesQ == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	currencies := strings.Split(currenciesQ, ",")
	if len(currencies) < 2 {
		c.Status(http.StatusBadRequest)
		return
	}

	rates, err := r.RatesService.GetRates(currencies)
	if err != nil {
		slog.Error("GetRates err", "currencies", currencies, "err", err)
		c.Status(http.StatusBadRequest)
		return
	}

	curMap := GetCurrenciesMap(currencies)

	c.JSON(http.StatusOK, gin.H{"currencies": curMap, "rates": rates})
}

func GetCurrenciesMap(currencies []string) map[string][]string {
	rates := make(map[string][]string, len(currencies))

	for i, c1 := range currencies {
		for j, c2 := range currencies {
			if i == j {
				continue
			}
			rates[c1] = append(rates[c1], c2)
		}
	}

	return rates
}
