package scraper

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type mirrorPayload struct {
	ID int    `json:"id"`
	I  int    `json:"i"`
	Q  string `json:"q"`
}

func decodeMirrorToken(token string) (*mirrorPayload, error) {
	if len(token) > 10240 {
		return nil, fmt.Errorf("token too large")
	}

	raw, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	var p mirrorPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *AnimeScraper) ResolveStreamURL(
	ctx context.Context,
	token string,
) (string, error) {

	nonce, err := s.getNonce(ctx)
	if err != nil {
		return "", err
	}

	payload, err := decodeMirrorToken(token)
	if err != nil {
		return "", err
	}

	form := url.Values{}
	form.Set("action", ActionGetEmbed)
	form.Set("nonce", nonce)
	form.Set("id", strconv.Itoa(payload.ID))
	form.Set("i", strconv.Itoa(payload.I))
	form.Set("q", payload.Q)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.baseURL+"/wp-admin/admin-ajax.php",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("create embed request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("embed request failed with status %d", resp.StatusCode)
	}

	var res struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if res.Data == "" {
		return "", fmt.Errorf("embed data empty")
	}

	if len(res.Data) > 5242880 { // 5MB limit
		return "", fmt.Errorf("embed data too large")
	}

	decoded, err := base64.StdEncoding.DecodeString(res.Data)
	if err != nil {
		return "", fmt.Errorf("decode embed data: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(decoded)))
	if err != nil {
		return "", fmt.Errorf("parse embed html: %w", err)
	}

	src, exists := doc.Find("iframe").Attr("src")
	if !exists {
		return "", fmt.Errorf("iframe src not found")
	}

	if !strings.HasPrefix(src, "http://") && !strings.HasPrefix(src, "https://") {
		return "", fmt.Errorf("invalid iframe src scheme")
	}

	return src, nil
}
