# Kunime API

Lightweight Go + Fiber service that scrapes anime data from [Otakudesu](https://otakudesu.best) (ongoing, completed, and by genre) and exposes it via REST. Every endpoint requires the `X-API-Key` header.

## Tech Stack

- Go (Fiber v2) for the HTTP server
- Colly for web scraping
- Simple logging and API key middleware

## Prerequisites

- Go 1.25+ (see [`go.mod`](go.mod))
- Internet access to the target `SCRAPE_BASE_URL`

## Environment Configuration

Create a `.env` file in the project root (or set environment variables directly):

```env
PORT=8080
API_KEY=supersecret
SCRAPE_BASE_URL=https://otakudesu.best
USER_AGENT=Mozilla/5.0 ...
```

| Variable          | Description                            | Default |
| ----------------- | -------------------------------------- | ------- |
| `API_KEY`         | API key required in the request header | -       |
| `SCRAPE_BASE_URL` | Base URL of the data source            | -       |
| `USER_AGENT`      | User-Agent string for Colly            | -       |
| `PORT`            | Fiber server port                      | `8080`  |

## Run Locally

```bash
go run ./cmd/server/main.go

# or inline env:
API_KEY=supersecret SCRAPE_BASE_URL=https://otakudesu.best USER_AGENT="Mozilla/5.0 ..." go run ./cmd/server/main.go
```

The server will be available at `http://localhost:<PORT>`

## Endpoints

All endpoints require the header `X-API-Key: <API_KEY>`

| Method | Path                                 | Description                                  |
| ------ | ------------------------------------ | -------------------------------------------- |
| GET    | `/`                                  | Quick info and endpoint list                 |
| GET    | `/healthz`                           | Health check                                 |
| GET    | `/api/v1/ongoing-anime/:page`        | Ongoing anime list (page defaults 1)         |
| GET    | `/api/v1/completed-anime/:page`      | Completed anime list (page defaults 1)       |
| GET    | `/api/v1/genres`                     | List of available genres                     |
| GET    | `/api/v1/genre/:genreSlug/:page`     | Anime by genre and page                      |
| GET    | `/api/v1/anime/:animeSlug/batch`     | Anime batch download links                   |
| GET    | `/api/v1/anime/:animeSlug`           | Anime detail                                 |
| GET    | `/api/v1/anime/:animeSlug/episodes`  | Anime episode list                           |
| GET    | `/api/v1/search/:query`              | Search anime                                 |
| GET    | `/api/v1/anime/:episodeSlug/streams` | List available streaming mirrors (tokenized) |
| POST   | `/api/v1/streams/resolve`            | Resolve stream token into final iframe URL   |

Examples:

```bash
curl -H "X-API-Key: supersecret" http://localhost:8080/api/v1/ongoing-anime/1
```

Example response (shortened):

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

## Documentation

Full API details live in [`docs/API.md`](docs/API.md)

## Project Structure (short)

- `cmd/server` – application entrypoint
- `internal/config` – env loader and configuration
- `internal/http` – Fiber router and handlers
- `internal/anime` – domain models and service
- `internal/scraper` – scraping logic with Colly
- `internal/middleware` – logging and API key middleware

## Buy me a coffee

If you find this project helpful, consider buying me a coffee!

<a href="https://buymeacoffee.com/kudanil" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>

## License

The project is licensed under the Apache License 2.0 - See [`LICENSE`](LICENSE) for more details
