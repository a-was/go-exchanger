package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type router struct{}

func RegisterRoutes(r *gin.Engine) {
	router := router{}

	r.GET("/rates", router.getRates)
}

func (r *router) getRates(c *gin.Context) {
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

	curMap := GetCurrenciesMap(currencies)

	c.JSON(http.StatusOK, gin.H{"currencies": curMap})
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
