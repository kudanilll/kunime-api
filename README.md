# Kunime API

Lightweight Go + Fiber service that scrapes anime data from [Otakudesu](https://otakudesu.best) (ongoing, completed, and by genre) and exposes it via REST. All API endpoints require the `X-API-Key` header except `GET /healthz`.

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

The server listens on `0.0.0.0:<PORT>` and will be available at `http://localhost:<PORT>`

## Endpoints

All endpoints require the header `X-API-Key: <API_KEY>` except `GET /healthz`

| Method | Path                                 | Description                            |
| ------ | ------------------------------------ | -------------------------------------- |
| GET    | `/`                                  | Quick info and endpoint list           |
| GET    | `/healthz`                           | Health check, no API key required      |
| GET    | `/api/v1/ongoing-anime/:page`        | Ongoing anime list (page defaults 1)   |
| GET    | `/api/v1/completed-anime/:page`      | Completed anime list (page defaults 1) |
| GET    | `/api/v1/genres`                     | List of available genres               |
| GET    | `/api/v1/genre/:genreSlug/:page`     | Anime by genre and page                |
| GET    | `/api/v1/anime/:animeSlug/batch`     | Anime batch download links             |
| GET    | `/api/v1/anime/:animeSlug`           | Anime detail                           |
| GET    | `/api/v1/anime/:animeSlug/episodes`  | Anime episode list                     |
| GET    | `/api/v1/anime/:episodeSlug/streams` | List available streaming mirrors       |
| GET    | `/api/v1/search/:query`              | Search anime                           |
| POST   | `/api/v1/streams/resolve`            | Resolve stream token into final url    |

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

## Support me

If you find this project helpful, consider buying me a coffee!

<a href="https://buymeacoffee.com/kudanil" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>

<div>
  <a href="https://saweria.co/achmaddaniel" target="_blank"><img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcR381toYEI4e_5hJ3nvA5pzj2PrtNc42FvBgQ&s" alt="Saweria" ></a>
</div>

## License

This project is licensed under the **GNU Affero General Public License v3.0**.

See the [`LICENSE`](LICENSE) file for full license text.

```
Copyright (C) 2026 Achmad Daniel Syahputra, Rizky Bintang Assabil

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
```
