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
		return "", fmt.Errorf("embed data empty")
	}

	decoded, _ := base64.StdEncoding.DecodeString(res.Data)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(decoded)))

	src, exists := doc.Find("iframe").Attr("src")
	if !exists {
		return "", fmt.Errorf("iframe src not found")
	}

	return src, nil
}
