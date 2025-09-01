# RSSHub: Агрегатор RSS-лент

## Описание

RSSHub - это CLI-приложение для агрегации и обработки RSS-лент. Система периодически получает статьи из пользовательских RSS-лент, сохраняет их в PostgreSQL и предоставляет интерфейс для управления подписками и просмотра новостей.

## Особенности

- **Фоновое обновление**: Периодическое получение новых статей с настраиваемым интервалом
- **Параллельная обработка**: Worker pool для конкурентной обработки RSS-лент
- **Управление через CLI**: Полный набор команд для управления подписками и настройками
- **Хранение данных**: PostgreSQL для надежного хранения лент и статей
- **Динамическая конфигурация**: Изменение интервала и количества воркеров без перезапуска

## Архитектура

Проект использует модульную структуру:

- **CLI Layer**: Обработка команд и взаимодействие с пользователем
- **Application Layer**: Бизнес-логика и use cases
- **Adapter Layer**: 
  - RSS Adapter (получение и парсинг RSS-лент)
  - Postgres Adapter (работа с базой данных)
- **Domain Layer**: Модели данных и интерфейсы

## Запуск проекта

### Предварительные требования

- Docker и Docker Compose
- Go 1.21+ (для локальной разработки)

### Настройка окружения

1. Скопируйте файл `.env.example` в `.env`:
```bash
cp .env.example .env
```

2. Заполните переменные окружения в файле `.env`:

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

### Запуск с Docker Compose

```bash
docker-compose up -d
```

### Настройка базы данных

После запуска PostgreSQL выполните миграции:

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

### Сборка и запуск приложения

```bash
go build -o rsshub .
./rsshub --help
```

## Использование

### Командная строка

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

### Команды

#### Добавление RSS-ленты

```bash
./rsshub add --name "tech-crunch" --url "https://techcrunch.com/feed/"
```

#### Запуск фонового процесса

```bash
./rsshub fetch
```

#### Изменение интервала обновления

```bash
./rsshub set-interval 2m
```

#### Изменение количества воркеров

```bash
./rsshub set-workers 5
```

#### Список RSS-лент

```bash
./rsshub list --num 5
```

#### Удаление RSS-ленты

```bash
./rsshub delete --name "tech-crunch"
```

#### Просмотр статей

```bash
./rsshub articles --feed-name "tech-crunch" --num 5
```

## Структура проекта

```
rsshub/
├── cmd/                 # Точка входа приложения
├── internal/            # Внутренние пакеты
│   ├── adapter/         # Адаптеры
│   │   ├── handlers/    # Обработчики
│   │   ├── postgres/    # PostgreSQL адаптер
│   │   └── rss/         # RSS адаптер
│   ├── app/             # Приложение
│   ├── apperrors/       # Ошибки приложения
│   ├── cli/             # CLI обработка
│   ├── config/          # Конфигурация
│   ├── domain/          # Доменный слой
│   ├── logger/          # Логирование
│   └── utils/           # Вспомогательные утилиты
├── migrations/          # Миграции базы данных
└── docker-compose.yml   # Docker Compose конфигурация
```

## База данных

### Таблица `feeds`

Хранит метаданные RSS-лент:

| Поле | Тип | Описание |
|------|-----|----------|
| id | UUID | Уникальный идентификатор |
| created_at | TIMESTAMP | Время создания |
| updated_at | TIMESTAMP | Время последнего обновления |
| name | TEXT | Человекочитаемое название |
| url | TEXT | URL RSS-ленты |

### Таблица `articles`

Хранит статьи из RSS-лент:

| Поле | Тип | Описание |
|------|-----|----------|
| id | UUID | Уникальный идентификатор |
| created_at | TIMESTAMP | Время создания |
| updated_at | TIMESTAMP | Время последнего обновления |
| title | TEXT | Заголовок статьи |
| link | TEXT | URL статьи |
| published_at | TIMESTAMP | Время публикации |
| description | TEXT | Описание статьи |
| feed_id | UUID | Ссылка на RSS-ленту |

## Примеры RSS-лент для тестирования

- TechCrunch: `https://techcrunch.com/feed/`
- Hacker News: `https://news.ycombinator.com/rss`
- UN News: `https://news.un.org/feed/subscribe/ru/news/all/rss.xml`
- BBC News: `https://feeds.bbci.co.uk/news/world/rss.xml`
- Ars Technica: `http://feeds.arstechnica.com/arstechnica/index`
- The Verge: `https://www.theverge.com/rss/index.xml`

## Graceful Shutdown

Приложение корректно обрабатывает сигналы завершения (Ctrl+C), останавливая все фоновые процессы и закрывая соединения с базой данных.