# EHIN API

This is the Go implementation of the EHIN API.

## Development

### Prerequisites

- Go 1.25 or later
- A PostgreSQL database (optional, for local development)

### Running locally

```bash
go run cmd/api/main.go
```

The server will start on port 8080 by default. You can configure it using environment variables or a `.env` file.

### Environment Variables

- `PORT`: Port to listen on (default: 8080)
- `DATABASE_URL`: PostgreSQL connection string.
- `UPDATE_PRICES_PASSWORD`: Password for the `/api/update-prices` endpoints
- `CORS_ALLOWED_ORIGINS`: Comma-separated list of allowed origins

## Database Migrations

Database migrations are handled by [Goose](https://github.com/pressly/goose). They are run as part of the release pipeline.

Migrations are located in `internal/migrations`.

To create a new migration (requires goose CLI):
```bash
goose -dir internal/migrations create your_migration_name sql
```

To run migrations locally:
```bash
goose -dir internal/migrations postgres "your_database_url" up
```

## Testing

Run all tests:

```bash
go test ./...
```

## Deployment

The API is configured for Google App Engine.

To build and deploy:

```bash
go build -o application cmd/api/main.go
gcloud app deploy
```

To deploy the cron jobs:

```bash
gcloud app deploy backup_cron.yaml
```