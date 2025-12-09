package config

import "os"

type Config struct {
    Port          string
    APIKey        string
    ScrapeBaseURL string
    UserAgent     string
}

func Load() Config {
    return Config{
        Port:          getEnv("PORT", "8080"),
        APIKey:        mustEnv("API_KEY"),
        ScrapeBaseURL: mustEnv("SCRAPE_BASE_URL"),
        UserAgent:     mustEnv("USER_AGENT"),
    }
}

func getEnv(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}

func mustEnv(key string) string {
    v := os.Getenv(key)
    if v == "" {
        panic("missing env: " + key)
    }
    return v
}
