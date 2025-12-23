# Kunime API Documentation

## Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
- [Authentication](#authentication)
- [Base URLs](#base-urls)
- [API Endpoints](#api-endpoints)
- [Streaming Flow](#streaming-flow)
- [Error Responses](#error-responses)

## Overview

Kunime API is a high-performance REST API built with Go and Fiber that provides comprehensive anime data scraped from Otakudesu. Access information about ongoing anime, completed series, genres, episode lists, download links, and more through clean, RESTful endpoints.

**Key Features:**

- Real-time anime data from Otakudesu
- Secure API key authentication
- Pagination support for large datasets
- Genre-based filtering
- Episode and download link information
- Search functionality

## Getting Started

### Prerequisites

- Go 1.25+ or higher
- Valid API key

### Environment Configuration

Create a `.env` file in your project root:

```env
API_KEY=your_secure_api_key_here
SCRAPE_BASE_URL=https://otakudesu.best
USER_AGENT=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36
PORT=8080
```

| Variable          | Required | Description                             | Default |
| ----------------- | -------- | --------------------------------------- | ------- |
| `API_KEY`         | Yes      | Secret key for API authentication       | -       |
| `SCRAPE_BASE_URL` | Yes      | Source website base URL                 | -       |
| `USER_AGENT`      | Yes      | User-Agent string for scraping requests | -       |
| `PORT`            | No       | Server port number                      | `8080`  |

### Installation & Running

```bash
# Clone the repository
git clone https://github.com/kudanilll/kunime-api.git
cd kunime-api

# Run with environment file
go run ./cmd/server/main.go

# Or run with inline environment variables
API_KEY=supersecret \
SCRAPE_BASE_URL=https://otakudesu.best \
USER_AGENT="Mozilla/5.0 ..." \
PORT=8080 \
go run ./cmd/server/main.go
```

The server will start at `http://localhost:8080` (or your configured port).

## Authentication

All API endpoints (except health check) require authentication via API key.

**Header Format:**

```
X-API-Key: your_api_key_here
```

**Authentication Errors:**

| Status | Response                       | Description                |
| ------ | ------------------------------ | -------------------------- |
| 401    | `{"error": "invalid api key"}` | Missing or invalid API key |

## Base URLs

**Local Development:**

```
http://localhost:8080
```

**Production:**

```
https://your-domain.com
```

### System Endpoints

#### Service Information

```http
GET /
```

Returns service information and available endpoints.

**Response:**

```json
{
  "endpoint": {
    "Get Anime Batch": "/api/v1/anime/:animeSlug/batch",
    "Get Anime Detail": "/api/v1/anime/:animeSlug",
    "Get Anime Episodes": "/api/v1/anime/:animeSlug/episodes",
    "Get Anime by Genre & Page": "/api/v1/genre/:genreSlug/:page",
    "Get Completed Anime": "/api/v1/completed-anime/:page",
    "Get Genres": "/api/v1/genres",
    "Get Ongoing Anime": "/api/v1/ongoing-anime/:page",
    "Search Anime": "/api/v1/search/:query"
  },
  "github": "https://github.com/kudanilll/kunime-api",
  "support": "https://buymeacoffee.com/kudanil"
}
```

#### Health Check

```http
GET /healthz
```

Returns service health status. **No authentication required.**

**Response:**

```json
{
  "status": "ok"
}
```

## API Endpoints

All endpoints require the `X-API-Key` header unless otherwise specified.

### 1. Ongoing Anime

Get list of currently airing anime series.

**Endpoints:**

```http
GET /api/v1/ongoing-anime
GET /api/v1/ongoing-anime/:page
```

**Parameters:**

| Name   | Type    | Location | Required | Description               |
| ------ | ------- | -------- | -------- | ------------------------- |
| `page` | integer | path     | No       | Page number (starts at 1) |

**Example Request:**

```bash
curl -X GET "http://localhost:8080/api/v1/ongoing-anime/1" \
  -H "X-API-Key: your_api_key"
```

**Success Response (200 OK):**

```json
{
  "page": 1,
  "data": [
    {
      "title": "One Piece",
      "episode": 1087,
      "day": "Sunday",
      "date": "22 Dec",
      "image": "https://otakudesu.best/wp-content/uploads/2024/one-piece.jpg",
      "endpoint": "https://otakudesu.best/anime/one-piece-sub-indo/"
    },
    {
      "title": "Dandadan",
      "episode": 12,
      "day": "Thursday",
      "date": "19 Dec",
      "image": "https://otakudesu.best/wp-content/uploads/2024/dandadan.jpg",
      "endpoint": "https://otakudesu.best/anime/dandadan-sub-indo/"
    }
  ]
}
```

### 2. Completed Anime

Get list of completed anime series.

**Endpoints:**

```http
GET /api/v1/completed-anime
GET /api/v1/completed-anime/:page
```

**Parameters:**

| Name   | Type    | Location | Required | Description               |
| ------ | ------- | -------- | -------- | ------------------------- |
| `page` | integer | path     | No       | Page number (starts at 1) |

**Example Request:**

```bash
curl -X GET "http://localhost:8080/api/v1/completed-anime/1" \
  -H "X-API-Key: your_api_key"
```

**Success Response (200 OK):**

```json
{
  "page": 1,
  "data": [
    {
      "title": "Frieren: Beyond Journey's End",
      "episodes": 28,
      "score": 9.37,
      "date": "22 Mar",
      "image": "https://otakudesu.best/wp-content/uploads/2023/frieren.jpg",
      "endpoint": "https://otakudesu.best/anime/frieren-sub-indo/"
    }
  ]
}
```

### 3. Genres

Get list of all available anime genres.

**Endpoint:**

```http
GET /api/v1/genres
```

**Example Request:**

```bash
curl -X GET "http://localhost:8080/api/v1/genres" \
  -H "X-API-Key: your_api_key"
```

**Success Response (200 OK):**

```json
{
  "data": [
    {
      "name": "Action",
      "slug": "action",
      "endpoint": "https://otakudesu.best/genre/action/"
    },
    {
      "name": "Adventure",
      "slug": "adventure",
      "endpoint": "https://otakudesu.best/genre/adventure/"
    },
    {
      "name": "Comedy",
      "slug": "comedy",
      "endpoint": "https://otakudesu.best/genre/comedy/"
    }
  ]
}
```

### 4. Anime by Genre

Get anime list filtered by specific genre with pagination.

**Endpoint:**

```http
GET /api/v1/genre/:genreSlug/:page
```

**Parameters:**

| Name        | Type    | Location | Required | Description                 |
| ----------- | ------- | -------- | -------- | --------------------------- |
| `genreSlug` | string  | path     | Yes      | Genre slug (e.g., "action") |
| `page`      | integer | path     | Yes      | Page number (starts at 1)   |

**Example Request:**

```bash
curl -X GET "http://localhost:8080/api/v1/genre/action/1" \
  -H "X-API-Key: your_api_key"
```

**Success Response (200 OK):**

```json
{
  "genre": "action",
  "page": 1,
  "data": [
    {
      "title": "Jujutsu Kaisen Season 2",
      "endpoint": "https://otakudesu.best/anime/jjk-s2-sub-indo/",
      "studio": "MAPPA",
      "episodes": "23 Eps",
      "rating": "8.91",
      "genres": ["Action", "Fantasy", "Shounen", "Supernatural"],
      "image": "https://otakudesu.best/wp-content/uploads/2023/jjk-s2.jpg",
      "synopsis": "Melanjutkan petualangan Yuji Itadori dalam dunia jujutsu...",
      "season": "Summer 2023"
    }
  ]
}
```

### 5. Anime Detail

Get detailed information about a specific anime.

**Endpoint:**

```http
GET /api/v1/anime/:animeSlug
```

**Parameters:**

| Name        | Type   | Location | Required | Description                                     |
| ----------- | ------ | -------- | -------- | ----------------------------------------------- |
| `animeSlug` | string | path     | Yes      | Anime slug identifier (e.g., "jjk-s2-sub-indo") |

**Example Request:**

```bash
curl -X GET "http://localhost:8080/api/v1/anime/kakkou-iinazuke-s2-sub-indo" \
  -H "X-API-Key: your_api_key"
```

**Success Response (200 OK):**

```json
{
  "title": "Kakkou no Iinazuke Season 2",
  "japanese_title": "カッコウの許嫁 Season2",
  "score": "6.72",
  "type": "TV",
  "status": "Completed",
  "total_episode": "12",
  "duration": "23 Menit",
  "release_date": "Jul 08, 2025",
  "studio": "Okuruto Noboru",
  "producers": ["Kodansha", "Crunchyroll", "BS NTV", "Kadokawa"],
  "genres": ["Comedy", "Harem", "Romance", "Shounen"],
  "image": "https://otakudesu.best/wp-content/uploads/2025/12/kakkou-s2.jpg",
  "synopsis": "Musim kedua dari anime Kakkou no Iinazuke yang menceritakan..."
}
```

### 6. Episode List

Get all available episodes for a specific anime.

**Endpoint:**

```http
GET /api/v1/anime/:animeSlug/episodes
```

**Parameters:**

| Name        | Type   | Location | Required | Description                                     |
| ----------- | ------ | -------- | -------- | ----------------------------------------------- |
| `animeSlug` | string | path     | Yes      | Anime slug identifier (e.g., "jjk-s2-sub-indo") |

**Example Request:**

```bash
curl -X GET "http://localhost:8080/api/v1/anime/kakkou-iinazuke-s2-sub-indo/episodes" \
  -H "X-API-Key: your_api_key"
```

**Success Response (200 OK):**

```json
{
  "anime_slug": "kakkou-iinazuke-s2-sub-indo",
  "episodes": [
    {
      "episode": 1,
      "slug": "kni-s2-episode-1-sub-indo"
    },
    {
      "episode": 2,
      "slug": "kni-s2-episode-2-sub-indo"
    },
    {
      "episode": 3,
      "slug": "kni-s2-episode-3-sub-indo"
    },
    {
      "episode": 4,
      "slug": "kni-s2-episode-4-sub-indo"
    },
    {
      "episode": 5,
      "slug": "kni-s2-episode-5-sub-indo"
    },
    {
      "episode": 6,
      "slug": "kni-s2-episode-6-sub-indo"
    },
    {
      "episode": 7,
      "slug": "kni-s2-episode-7-sub-indo"
    },
    {
      "episode": 8,
      "slug": "kni-s2-episode-8-sub-indo"
    },
    {
      "episode": 9,
      "slug": "kni-s2-episode-9-sub-indo"
    },
    {
      "episode": 10,
      "slug": "kni-s2-episode-10-sub-indo"
    },
    {
      "episode": 11,
      "slug": "kni-s2-episode-11-sub-indo"
    },
    {
      "episode": 12,
      "slug": "kni-s2-episode-12-sub-indo"
    }
  ]
}
```

### 7. Batch Download Links

Get batch download links for all episodes of an anime in various qualities.

**Endpoint:**

```http
GET /api/v1/anime/:animeSlug/batch
```

**Parameters:**

| Name        | Type   | Location | Required | Description                                      |
| ----------- | ------ | -------- | -------- | ------------------------------------------------ |
| `animeSlug` | string | path     | Yes      | Anime batch slug (e.g., "kni-s2-batch-sub-indo") |

**Example Request:**

```bash
curl -X GET "http://localhost:8080/api/v1/anime/kni-s2-batch-sub-indo/batch" \
  -H "X-API-Key: your_api_key"
```

**Success Response (200 OK):**

```json
{
  "title": "Kakkou no Iinazuke Season 2 Batch Subtitle Indonesia",
  "qualities": [
    {
      "quality": "MP4 360p",
      "size": "0.46 GB",
      "links": [
        {
          "server": "OtakuDrive",
          "url": "https://desustream.com/safelink/..."
        },
        {
          "server": "DesuDrive",
          "url": "https://desustream.com/safelink/..."
        }
      ]
    },
    {
      "quality": "MP4 480p",
      "size": "0.79 GB",
      "links": [
        {
          "server": "Mega",
          "url": "https://desustream.com/safelink/..."
        },
        {
          "server": "GDrive",
          "url": "https://desustream.com/safelink/..."
        }
      ]
    },
    {
      "quality": "MP4 720p",
      "size": "1.52 GB",
      "links": [
        {
          "server": "AceFile",
          "url": "https://desustream.com/safelink/..."
        }
      ]
    },
    {
      "quality": "MP4 1080p",
      "size": "2.84 GB",
      "links": [
        {
          "server": "Pixeldrain",
          "url": "https://desustream.com/safelink/..."
        }
      ]
    }
  ]
}
```

### 8. Search Anime

Search for anime by title or keywords.

**Endpoint:**

```http
GET /api/v1/search/:query
```

**Parameters:**

| Name    | Type   | Location | Required | Description                          |
| ------- | ------ | -------- | -------- | ------------------------------------ |
| `query` | string | path     | Yes      | Search keywords (use `+` for spaces) |

**Example Request:**

```bash
curl -X GET "http://localhost:8080/api/v1/search/jujutsu+kaisen" \
  -H "X-API-Key: your_api_key"
```

**Success Response (200 OK):**

```json
{
  "query": "jujutsu+kaisen",
  "data": [
    {
      "title": "Jujutsu Kaisen Season 2 Sub Indo",
      "status": "Completed",
      "rating": "8.91",
      "genres": ["Action", "Fantasy", "Shounen", "Supernatural"],
      "image": "https://otakudesu.best/wp-content/uploads/2023/jjk-s2.jpg",
      "endpoint": "https://otakudesu.best/anime/jjk-s2-sub-indo/"
    },
    {
      "title": "Jujutsu Kaisen 0 Movie Sub Indo",
      "status": "Completed",
      "rating": "8.79",
      "genres": ["Action", "Fantasy", "Shounen", "Supernatural"],
      "image": "https://otakudesu.best/wp-content/uploads/2022/jjk-0.jpg",
      "endpoint": "https://otakudesu.best/anime/jjk-0-movie-sub-indo/"
    }
  ]
}
```

### 9. Episode Streaming Mirrors

Get available streaming mirrors for a specific episode.  
Response contains **quality**, **server**, and **token**.  
The token must be resolved to obtain the final streaming URL.

**Endpoint:**

```http
GET /api/v1/anime/:episodeSlug/streams
```

**Parameters:**

| Name          | Type   | Location | Required | Description                                     |
| ------------- | ------ | -------- | -------- | ----------------------------------------------- |
| `episodeSlug` | string | path     | Yes      | Episode slug (e.g. `kni-s2-episode-1-sub-indo`) |

**Example Request:**

```bash
curl -X GET \
http://localhost:8080/api/v1/anime/kni-s2-episode-1-sub-indo/streams \
-H "X-API-Key: your_api_key"
```

**Success Response (200 OK):**

```json
{
  "episode_slug": "kni-s2-episode-1-sub-indo",
  "streams": [
    {
      "quality": "480p",
      "server": "otakuwatch5",
      "token": "eyJpZCI6MTkwNDE1LCJpIjowLCJxIjoiNDgwcCJ9"
    },
    {
      "quality": "720p",
      "server": "mega",
      "token": "eyJpZCI6MTkwNDE1LCJpIjoyLCJxIjoiNzIwcCJ9"
    }
  ]
}
```

### 10. Resolve Stream URL

Resolve a stream token into the final iframe streaming URL.

This endpoint performs a backend WordPress AJAX request and returns a direct iframe `src` URL.
Recommended usage: client selects mirror → backend resolves token.

**Endpoint:**

```http
POST /api/v1/streams/resolve
```

**Request Body:**

```json
{
  "token": "base64-encoded-token"
}
```

**Example Request:**

```bash
curl -X POST http://localhost:8080/api/v1/streams/resolve \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{"token":"eyJpZCI6MTkwNDE1LCJpIjowLCJxIjoiNDgwcCJ9"}'
```

**Success Response (200 OK):**

```json
{
  "url": "https://desustream.info/dstream/otakuwatch5/v3/index.php?id=VzcwdUErZXJBUVhQ..."
}
```

### Streaming Flow

1. Get episode streams (mirrors + tokens)

```bash
curl -H "X-API-Key: supersecret" \
http://localhost:8080/api/v1/anime/kni-s2-episode-1-sub-indo/streams
```

**Error Responses**

The API uses standard HTTP status codes to indicate success or failure.

```json
{
  "error": "Error message description"
}
```

2. Resolve a selected stream token into final streaming URL

```bash
curl -X POST http://localhost:8080/api/v1/streams/resolve \
  -H "X-API-Key: supersecret" \
  -H "Content-Type: application/json" \
  -d '{"token":"<base64-token>"}'
```

3. Streaming Resolution

   - Streaming URLs are not static.
   - Each mirror and resolution must be resolved individually using `/streams/resolve`.
   - Tokens may expire and should not be cached long-term.

### Common Status Codes

| Code | Status                | Description                             |
| ---- | --------------------- | --------------------------------------- |
| 200  | OK                    | Request succeeded                       |
| 400  | Bad Request           | Invalid parameters or malformed request |
| 401  | Unauthorized          | Missing or invalid API key              |
| 404  | Not Found             | Resource not found                      |
| 500  | Internal Server Error | Server error or scraping failure        |

### Error Examples

**401 Unauthorized:**

```json
{
  "error": "invalid api key"
}
```

**400 Bad Request:**

```json
{
  "error": "invalid page parameter"
}
```

**500 Internal Server Error:**

```json
{
  "error": "failed to scrape data from source"
}
```

## Notes & Limitations

**Important Considerations:**

1. **Data Source Dependency**: This API scrapes data from Otakudesu. Any changes to their website structure may temporarily affect API functionality.

2. **Pagination**:

   - Page numbers start at 1
   - Invalid page numbers default to 1
   - Pagination mirrors the source site structure

3. **Performance**: Response times depend on the source website's availability and response time.

4. **Slug Formats**:

   - Anime slugs typically end with `-sub-indo`
   - Batch slugs end with `-batch-sub-indo`
   - Use the search endpoint to find correct slugs

5. **Data Freshness**: Data is scraped in real-time, ensuring up-to-date information but potentially longer response times.

6. **Download Links**: Links in the batch endpoint may redirect through the source site's safelink system.

## Support & Contribution

For issues, feature requests, or contributions:

- Open an issue on GitHub
- Submit pull requests
- Contact: [hello.achmaddaniel@gmail.com](mailto:hello.achmaddaniel@gmail.com)

**License**: [Apache License 2.0](../LICENSE)

**Last Updated**: December 2025
