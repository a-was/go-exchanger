package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type OpenExchangeRatesService struct {
	AppID string
}

var _ RatesGetter = (*OpenExchangeRatesService)(nil)

type OpenExchangeResponse struct {
	Timestamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Rates     map[string]float64 `json:"rates"`
}

func (s *OpenExchangeRatesService) GetRates(targetCurrencies []string) (RatesMap, error) {
	url := fmt.Sprintf(
		"https://openexchangerates.org/api/latest.json?app_id=%s&symbols=%s",
		s.AppID, strings.Join(targetCurrencies, ","),
	)

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("http.Get err", "url", url, "err", err)
		return nil, fmt.Errorf("http.Get err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("http.Get invalid status code", "url", url, "status_code", resp.StatusCode)
		return nil, fmt.Errorf("http.Get invalid status code: %d", resp.StatusCode)
	}

	var response OpenExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.Error("json.NewDecoder err", "url", url, "err", err)
		return nil, fmt.Errorf("json.NewDecoder err: %w", err)
	}

	return buildRatesMap(response.Rates), nil
}
