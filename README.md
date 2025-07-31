# GitLab Readiness API

A Go API that tracks project production readiness by validating security and compliance checks.

## Requirements

- Go 1.21+
- PostgreSQL (database)
- Docker (optional, for local PostgreSQL)

## Quick Start

### Option 1: With Docker (Recommended)

```bash
# Clone and setup
git clone <repository-url>
cd go-backend

# Start PostgreSQL
docker-compose up -d postgres

# Run the application
go run cmd/api/main.go
```

### Option 2: With External PostgreSQL

```bash
# Set database connection
export DATABASE_URL="postgres://username:password@host:port/database?sslmode=disable"

# Run the application
go run cmd/api/main.go
```

### Option 3: Using Dev Containers

1. Open in VSCode
2. Click "Reopen in Container" when prompted
3. Run: `go run cmd/api/main.go`

The API will be available at `http://localhost:8080`

## Database Configuration

- **Local Development**: Uses Docker Compose PostgreSQL by default
- **Production**: Set `DATABASE_URL` environment variable
- **Default Local**: `postgres://postgres:password@localhost:5432/gitlab_readiness?sslmode=disable`

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check endpoint |
| GET | `/api/v1/gitlab/projects` | List all GitLab projects |
| GET | `/api/v1/gitlab/projects/{id}` | Get a single GitLab project |
| POST | `/api/v1/gitlab/projects` | Create a new GitLab project |
| PUT | `/api/v1/gitlab/projects/{id}` | Update an existing GitLab project |
| DELETE | `/api/v1/gitlab/projects/{id}` | Delete a GitLab project |

## API Documentation

Swagger UI: `http://localhost:8080/swagger/index.html`

Generate/update docs:
```bash
swag init -g cmd/api/main.go
```

## Project Structure

```
go-backend/
├── cmd/api/           # Application entry point
├── internal/          # Private application code
│   ├── config/        # Configuration management
│   ├── database/      # Database connection and migrations
│   ├── handlers/      # HTTP handlers
│   ├── models/        # Domain models
│   ├── repository/    # Data access layer
│   └── router/        # HTTP routing
├── migrations/        # SQL migration files
├── docs/              # Documentation
└── .devcontainer/     # Dev container configuration
```

## Configuration

The application uses environment variables for configuration. See `.env.example` for all available options.

Key configuration variables:
- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: `debug`, `info`, `warn`, or `error`

## Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...
```

## Building

```bash
go build -o bin/api cmd/api/main.go
```

## Database Schema

The application tracks the following GitLab readiness checks:

- Project presence in GitLab
- APP_NAME and MOAB_ID variables in .gitlab-ci.yml
- CODEOWNERS file existence
- Branch protection settings
- Merge request approval settings
- Commit message push rules

See [migrations/](migrations/) for the complete schema.

## Debugging in VSCode

1. Set breakpoints in your code
2. Press F5 or go to Run > Start Debugging
3. Select configuration:
   - **Launch API**: Debug with PostgreSQL
   - **Debug Current Test**: Debug the current test file

## Production

1. Configure `DATABASE_URL` with PostgreSQL connection string
2. Set `ENVIRONMENT=production`
3. Set `LOG_LEVEL=info` or `warn`