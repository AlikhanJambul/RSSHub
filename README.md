# RSSHub: RSS Feed Aggregator

## Description

**RSSHub** is a CLI application for aggregating and processing RSS feeds.
The system periodically fetches articles from user-defined RSS feeds, stores them in PostgreSQL, and provides a command-line interface for managing subscriptions and viewing news.

## Features

* **Background updates** — periodic fetching of new articles with a configurable interval
* **Parallel processing** — worker pool for concurrent RSS feed processing
* **CLI management** — full set of commands for managing feeds and settings
* **Data persistence** — PostgreSQL for reliable storage of feeds and articles
* **Dynamic configuration** — change update intervals and worker counts without restarting

## Architecture

The project follows a modular architecture:

* **CLI Layer** — handles commands and user interaction
* **Application Layer** — business logic and use cases
* **Adapter Layer**:

  * **RSS Adapter** — fetching and parsing RSS feeds
  * **Postgres Adapter** — database interaction
* **Domain Layer** — data models and interfaces

## Running the Project

### Prerequisites

* Docker and Docker Compose
* Go 1.21+ (for local development)

### Environment Configuration

1. Copy the `.env.example` file to `.env`:

```bash
cp .env.example .env
```

2. Fill in the environment variables in `.env`:

```env
# CLI App
CLI_APP_TIMER_INTERVAL=3m
CLI_APP_WORKERS_COUNT=3

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=changeme
POSTGRES_DBNAME=rsshub
```

### Running with Docker Compose

```bash
docker-compose up -d
```

### Database Setup

After PostgreSQL starts, run the migrations:

```bash
docker exec -it rsshub_postgres psql -U postgres -d rsshub -c "
CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";

CREATE TABLE feeds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    name TEXT UNIQUE NOT NULL,
    url TEXT NOT NULL
);

CREATE TABLE articles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    title TEXT NOT NULL,
    link TEXT UNIQUE NOT NULL,
    published_at TIMESTAMP,
    description TEXT,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX articles_feed_link_idx ON articles(feed_id, link);
"
```

### Building and Running the Application

```bash
go build -o rsshub .
./rsshub --help
```

## Usage

### Command Line

```bash
./rsshub --help

Usage:
  rsshub COMMAND [OPTIONS]

Common Commands:
  add             add new RSS feed
  set-interval    set RSS fetch interval
  set-workers     set number of workers
  list            list available RSS feeds
  delete          delete RSS feed
  articles        show latest articles
  fetch           starts the background process that periodically fetches and processes RSS feeds using a worker pool
```

### Commands

#### Add an RSS Feed

```bash
./rsshub add --name "tech-crunch" --url "https://techcrunch.com/feed/"
```

#### Start Background Fetching

```bash
./rsshub fetch
```

#### Change Update Interval

```bash
./rsshub set-interval 2m
```

#### Change Number of Workers

```bash
./rsshub set-workers 5
```

#### List RSS Feeds

```bash
./rsshub list --num 5
```

#### Delete an RSS Feed

```bash
./rsshub delete --name "tech-crunch"
```

#### View Articles

```bash
./rsshub articles --feed-name "tech-crunch" --num 5
```

## Project Structure

```
rsshub/
├── cmd/                 # Application entry point
├── internal/            # Internal packages
│   ├── adapter/         # Adapters
│   │   ├── handlers/    # Handlers
│   │   ├── postgres/    # PostgreSQL adapter
│   │   └── rss/         # RSS adapter
│   ├── app/             # Application logic
│   ├── apperrors/       # Application errors
│   ├── cli/             # CLI handling
│   ├── config/          # Configuration
│   ├── domain/          # Domain layer
│   ├── logger/          # Logging
│   └── utils/           # Utility functions
├── migrations/          # Database migrations
└── docker-compose.yml   # Docker Compose configuration
```

## Database

### `feeds` Table

Stores metadata about RSS feeds:

| Field      | Type      | Description              |
| ---------- | --------- | ------------------------ |
| id         | UUID      | Unique identifier        |
| created_at | TIMESTAMP | Creation time            |
| updated_at | TIMESTAMP | Last update time         |
| name       | TEXT      | Human-readable feed name |
| url        | TEXT      | Feed URL                 |

### `articles` Table

Stores articles fetched from RSS feeds:

| Field        | Type      | Description               |
| ------------ | --------- | ------------------------- |
| id           | UUID      | Unique identifier         |
| created_at   | TIMESTAMP | Creation time             |
| updated_at   | TIMESTAMP | Last update time          |
| title        | TEXT      | Article title             |
| link         | TEXT      | Article URL               |
| published_at | TIMESTAMP | Publication date          |
| description  | TEXT      | Article description       |
| feed_id      | UUID      | Reference to the RSS feed |

## Example RSS Feeds for Testing

* TechCrunch: `https://techcrunch.com/feed/`
* Hacker News: `https://news.ycombinator.com/rss`
* UN News: `https://news.un.org/feed/subscribe/ru/news/all/rss.xml`
* BBC News: `https://feeds.bbci.co.uk/news/world/rss.xml`
* Ars Technica: `http://feeds.arstechnica.com/arstechnica/index`
* The Verge: `https://www.theverge.com/rss/index.xml`

## Graceful Shutdown

The application gracefully handles termination signals (Ctrl+C), stopping all background processes and closing database connections.
