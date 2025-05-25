package routes

import (
	"log/slog"
	"math/big"
	"net/http"
	"strconv"
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
	r.Engine.GET("/exchange", r.getExchange)
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

type cryptoRate struct {
	DecimalPlaces int
	Rate          float64
}

var db = map[string]cryptoRate{
	"BEER":  {DecimalPlaces: 18, Rate: 0.000_02461},
	"FLOKI": {DecimalPlaces: 18, Rate: 0.000_1428},
	"GATE":  {DecimalPlaces: 18, Rate: 6.87},
	"USDT":  {DecimalPlaces: 6, Rate: 0.999},
	"WBTC":  {DecimalPlaces: 8, Rate: 57037.22},
}

type preciseAmount struct {
	f    *big.Float
	prec int
}

func (a *preciseAmount) MarshalJSON() ([]byte, error) {
	if a == nil {
		return []byte("null"), nil
	}
	str := a.f.Text('f', a.prec)
	return []byte(str), nil
}

type outExchange struct {
	From   string         `json:"from"`
	To     string         `json:"to"`
	Amount *preciseAmount `json:"amount"`
}

func (r *Router) getExchange(c *gin.Context) {
	from := c.Query("from")
	if from == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	fromCryptoRate, ok := db[from]
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	to := c.Query("to")
	if to == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	toCryptoRate, ok := db[to]
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	amountQ := c.Query("amount")
	if amountQ == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	amount, err := strconv.ParseFloat(amountQ, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	const prec uint = 256

	scale := new(big.Float).SetPrec(prec).SetFloat64(1)
	ten := big.NewFloat(10)
	for range toCryptoRate.DecimalPlaces {
		scale.Mul(scale, ten)
	}

	fromRate := new(big.Float).SetPrec(prec).SetFloat64(fromCryptoRate.Rate)
	fromRate.Mul(fromRate, scale)
	toRate := new(big.Float).SetPrec(prec).SetFloat64(toCryptoRate.Rate)
	toRate.Mul(toRate, scale)

	amountBig := big.NewFloat(amount)
	amountInBase := new(big.Float).SetPrec(prec).Mul(amountBig, fromRate)
	amountInTarget := new(big.Float).SetPrec(prec).Quo(amountInBase, toRate)

	c.JSON(http.StatusOK, outExchange{
		From:   from,
		To:     to,
		Amount: &preciseAmount{amountInTarget, toCryptoRate.DecimalPlaces},
	})
}
