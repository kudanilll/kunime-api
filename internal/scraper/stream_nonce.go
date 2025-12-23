package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	ActionGetNonce = "aa1208d27f29ca340c92c66d1926f13f"
	ActionGetEmbed = "2a3505c93b0035d3f455df82bf976b84"
)

func (s *AnimeScraper) getNonce(ctx context.Context) (string, error) {
	form := url.Values{}
	form.Set("action", ActionGetNonce)

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.baseURL+"/wp-admin/admin-ajax.php",
		strings.NewReader(form.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if res.Data == "" {
		return "", fmt.Errorf("nonce empty")
	}

	return res.Data, nil
}
