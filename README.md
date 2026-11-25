# Kunime API

Kunime API supplies anime listings, search, and episode data to the [Kunime app](https://github.com/kudanilll/kunime). It fetches catalog pages from Otakudesu with Colly and serves JSON through Fiber. [Kunime Core](https://github.com/kudanilll/kunime-core) uses the completed-anime feed to build recommendations.

## Run Locally

Install Go 1.27.4 or newer, as specified in [`go.mod`](go.mod). Copy [`.env.example`](.env.example) to `.env` and fill in the required values before starting the server.

| Variable | What to set |
| --- | --- |
| `API_KEY` | The key clients send in the `X-API-Key` header. Required. |
| `SCRAPE_BASE_URL` | The upstream site URL. The template uses `https://otakudesu.best`; the server requires a value. |
| `USER_AGENT` | The user-agent string sent to the upstream site. Required. |
| `PORT` | HTTP port. Defaults to `8080`. |

```sh
go run ./cmd/server
```

In a second terminal, check that the server is listening:

```sh
curl --fail http://localhost:8080/healthz
```

Expect `{"status":"ok"}`. The health check only confirms the HTTP server is responding. To check upstream access, set `API_KEY` in your shell to the value used by the server and request a catalog page:

```sh
curl --fail -H "X-API-Key: $API_KEY" \
  http://localhost:8080/api/v1/ongoing-anime/1
```

## Endpoints

`GET /healthz` is public. Every other route, including `/`, requires `X-API-Key`. See the [API reference](docs/API.md) for request and response examples.

| Method | Path                                 | Description                            |
| ------ | ------------------------------------ | -------------------------------------- |
| GET    | `/`                                  | Quick info and endpoint list           |
| GET    | `/healthz`                           | Health check, no API key required      |
| GET    | `/api/v1/ongoing-anime[/:page]`       | Ongoing anime; defaults to page 1      |
| GET    | `/api/v1/completed-anime[/:page]`     | Completed anime; defaults to page 1    |
| GET    | `/api/v1/genres`                     | List of available genres               |
| GET    | `/api/v1/genre/:genreSlug/:page`     | Anime by genre and page                |
| GET    | `/api/v1/anime/:animeSlug/batch`     | Anime batch download links             |
| GET    | `/api/v1/anime/:animeSlug`           | Anime detail                           |
| GET    | `/api/v1/anime/:animeSlug/episodes`  | Anime episode list                     |
| GET    | `/api/v1/scrape-base-url`             | Configured source origin                |
| GET    | `/api/v1/anime/:episodeSlug/streams` | List available streaming mirrors       |
| GET    | `/api/v1/search/:query`              | Search anime                           |
| POST   | `/api/v1/streams/resolve`             | Resolve a mirror token to a stream URL |

## Tests

```sh
go test -mod=readonly ./...
```

The route tests cover the public health check, API-key handling, and invalid stream requests without contacting the upstream site.

## Deployment to Dokploy

The [GitHub Actions workflow](.github/workflows/docker-build-push.yml) runs Go checks and a Docker build on pull requests to `master`. A push to `master` runs the Go checks, pushes a Linux/amd64 image to GitHub Container Registry (GHCR) with `latest` and `sha-<full-commit>` tags, then calls Dokploy. [Dependabot](.github/dependabot.yml) checks Go modules, Actions, and base images weekly.

1. In Dokploy, select the **Docker** source. Set the image to `ghcr.io/kudanilll/kunime-api:latest` and point the domain to container port `8080`. If the package is private, configure GHCR credentials with `read:packages` and access to the package.
2. Set `API_KEY`, `SCRAPE_BASE_URL`, `USER_AGENT`, and `PORT=8080` in the application's environment.
3. Enable **Auto Deploy**. Copy the HTTPS webhook URL from the application's **Deployments** tab into the GitHub Actions secret `DOKPLOY_WEBHOOK_URL`, at repository level or in the `production` environment. Use one deployment trigger per image to avoid duplicate runs.
4. After a push to `master`, confirm the rollout in Dokploy. Check `GET /healthz` on the deployed domain, then make an authenticated catalog request. The webhook response only confirms Dokploy accepted the request.

To roll back, point Dokploy to a known-good `sha-<full-commit>` tag and deploy it manually. Switch back to `latest` when you want automatic deployment to resume.

## Project Structure

- `cmd/server` – application entrypoint
- `internal/config` – env loader and configuration
- `internal/http` – Fiber router and handlers
- `internal/anime` – domain models and service
- `internal/scraper` – scraping logic with Colly
- `internal/middleware` – logging and API key middleware

## Support

[Buy Me a Coffee](https://buymeacoffee.com/kudanil) · [Saweria](https://saweria.co/achmaddaniel)

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
