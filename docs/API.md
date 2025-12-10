# Kunime API - Documentation

## Overview

Kunime API is a Go + Fiber service that scrapes anime data from Otakudesu and exposes it through REST endpoints. Data includes ongoing anime, completed anime, available genres, and anime filtered by genre. All endpoints require an `X-API-Key` header for access.

## Authentication

- Header: `X-API-Key: <API_KEY>`
- If the header is missing or invalid, the server responds with `401 {"error": "invalid api key"}`.

## Environment Variables

| Name              | Description                                   | Default |
| ----------------- | --------------------------------------------- | ------- |
| `API_KEY`         | API key required on every request header      | -       |
| `SCRAPE_BASE_URL` | Base URL of the source site (e.g. Otakudesu)  | -       |
| `USER_AGENT`      | User-Agent string used by Colly when scraping | -       |
| `PORT`            | Port for the Fiber server                     | `8080`  |

## Running Locally

```bash
# with .env
go run ./cmd/server

# or inline
API_KEY=supersecret SCRAPE_BASE_URL=https://otakudesu.best USER_AGENT="Mozilla/5.0 ..." PORT=8080 go run ./cmd/server
```

Service runs on `http://localhost:<PORT>`.

## Base Routes

- `GET /` – service info and endpoint list
- `GET /healthz` – health check, returns `{"status": "ok"}`

## API Endpoints

All endpoints below require `X-API-Key`.

### 1) List Ongoing Anime

- `GET /api/v1/ongoing-anime`
- `GET /api/v1/ongoing-anime/:page`
- `page` (path, optional): starts at 1. If omitted or invalid, defaults to 1.

**Response**

```json
{
  "page": 1,
  "data": [
    {
      "title": "Anime Title",
      "episode": 10,
      "day": "Saturday",
      "date": "06 Dec",
      "image": "https://otakudesu.best/path/to/image.jpg",
      "endpoint": "https://otakudesu.best/anime/slug/"
    }
  ]
}
```

### 2) List Completed Anime

- `GET /api/v1/completed-anime`
- `GET /api/v1/completed-anime/:page`
- `page` (path, optional): starts at 1. If omitted or invalid, defaults to 1.

**Response**

```json
{
  "page": 1,
  "data": [
    {
      "title": "Anime Title",
      "episodes": 12,
      "score": 7.07,
      "date": "06 Dec",
      "image": "https://otakudesu.best/path/to/image.jpg",
      "endpoint": "https://otakudesu.best/anime/slug/"
    }
  ]
}
```

### 3) List Genres

- `GET /api/v1/genres`

**Response**

```json
{
  "data": [
    {
      "name": "Action",
      "slug": "action",
      "endpoint": "https://otakudesu.best/genre/action/"
    }
  ]
}
```

### 4) Anime by Genre & Page

- `GET /api/v1/genre/:genreSlug/:page`
- `genreSlug` (path, required): e.g. `action`
- `page` (path, required): starts at 1. If invalid, defaults to 1.

**Response**

```json
{
  "genre": "action",
  "page": 1,
  "data": [
    {
      "title": "Anime Title",
      "endpoint": "https://otakudesu.best/anime/slug/",
      "studio": "Studio Name",
      "episodes": "12 Eps",
      "rating": "7.20",
      "genres": ["Action", "Adventure"],
      "image": "https://otakudesu.best/path/to/image.jpg",
      "synopsis": "Short synopsis...",
      "season": "Fall 2025"
    }
  ]
}
```

## Error Handling

- `400` – invalid or missing path parameters
- `401` – missing/invalid API key
- `500` – upstream scrape errors or internal failures

## Notes & Limits

- The API scrapes live pages; availability and structure of the source site may affect responses.
- Pagination mirrors the source site: page `1` typically maps to the base listing, and `2+` map to paged URLs.
- Be respectful of the source site. The scraper uses a small random delay; consider adding rate limiting if you expose this publicly.
