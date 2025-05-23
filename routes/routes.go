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

type outRate struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
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

	response := make([]outRate, 0, len(rates)*(len(rates)-1))

	for from, rates := range rates {
		for to, rate := range rates {
			response = append(response, outRate{
				From: from,
				To:   to,
				Rate: rate,
			})
		}
	}

	c.JSON(http.StatusOK, response)
}
