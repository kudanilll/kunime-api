package scraper

import (
	"net/url"
	"path"
	"strconv"
	"strings"
)

// "Episode 10" -> 10
func extractEpisodeNumber(epText string) int {
	epText = strings.ToLower(epText)
	epText = strings.ReplaceAll(epText, "episode", "")
	epText = strings.TrimSpace(epText)

	if epText == "" {
		return 0
	}

	n, err := strconv.Atoi(epText)
	if err != nil {
		return 0
	}
	return n
}

func absoluteURL(base, p string) string {
	if p == "" {
		return ""
	}
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}

	u, err := url.Parse(base)
	if err != nil {
		return p
	}

	u.Path = path.Join(u.Path, p)
	return u.String()
}

func extractScore(text string) float64 {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}

	parts := strings.Fields(text)
	last := parts[len(parts)-1]

	f, err := strconv.ParseFloat(last, 64)
	if err != nil {
		return 0
	}
	return f
}

func extractGenreSlug(href string) string {
	href = strings.TrimSpace(href)
	if href == "" {
		return ""
	}

	// kalau absolute URL → ambil path-nya
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		u, err := url.Parse(href)
		if err != nil {
			return ""
		}
		href = u.Path
	}

	// "/genres/action/" → "genres/action"
	href = strings.Trim(href, "/")
	if href == "" {
		return ""
	}

	parts := strings.Split(href, "/")
	return parts[len(parts)-1]
}

func extractValue(text string) string {
	parts := strings.SplitN(text, ":", 2)
	if len(parts) != 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func extractEpisodeFromTitle(title string) int {
	title = strings.ToLower(title)

	// cari kata "episode"
	idx := strings.Index(title, "episode")
	if idx == -1 {
		return 0
	}

	part := title[idx+len("episode"):]
	part = strings.TrimSpace(part)

	fields := strings.Fields(part)
	if len(fields) == 0 {
		return 0
	}

	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0
	}

	return n
}
