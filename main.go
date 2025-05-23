package main

import (
	"log/slog"
	"os"

	"github.com/a-was/go-exchanger/routes"
	"github.com/a-was/go-exchanger/services"
	"github.com/gin-gonic/gin"
)

func main() {
	appID := os.Getenv("OPEN_EXCHANGE_RATES_APP_ID")
	if appID == "" {
		slog.Error("env OPEN_EXCHANGE_RATES_APP_ID not set")
		return
	}

	router := routes.Router{
		Engine: gin.Default(),
		RatesService: &services.OpenExchangeRatesService{
			AppID: appID,
		},
	}
	router.RegisterRoutes()

	slog.Info("Starting server", "port", "8080")
	router.Engine.Run(":8080")
}
