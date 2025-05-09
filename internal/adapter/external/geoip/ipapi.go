package geoip

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ip_detector/internal/logger"
)

type IPAPIService struct {
	APIURL string
}

type ipAPIResponse struct {
	Country string `json:"country"`
}

func NewIPAPIService(apiURL string) *IPAPIService {
	return &IPAPIService{APIURL: apiURL}
}

func (s *IPAPIService) GetCountryByIP(ip string) (string, error) {
	url := fmt.Sprintf("%s/%s?fields=country", s.APIURL, ip)
	logger.Log.Sugar().Infow("requesting GeoIP", "url", url, "ip", ip)

	resp, err := http.Get(url)
	if err != nil {
		logger.Log.Sugar().Errorw("failed to call GeoIP service", "error", err)
		return "", fmt.Errorf("failed to request IP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Log.Sugar().Warnw("GeoIP returned non‑200", "status", resp.Status, "ip", ip)
		return "", fmt.Errorf("IP API returned non-200 status: %s", resp.Status)
	}

	var data ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Log.Sugar().Errorw("failed to decode GeoIP response", "error", err)
		return "", fmt.Errorf("failed to decode IP API response: %w", err)
	}

	logger.Log.Sugar().Infow("GeoIP success", "ip", ip, "country", data.Country)
	return data.Country, nil
}
