# GitLab Readiness API

A Go-based backend API that tracks GitLab project production readiness by validating various security and compliance checks.

## Overview

This API provides CRUD operations for managing GitLab project readiness checks. It's designed as an educational project demonstrating Go best practices, clean architecture, and production-ready patterns.

## Features

- RESTful API for GitLab project readiness tracking
- Support for both SQLite (development) and PostgreSQL (production)
- Clean architecture with repository pattern
- Comprehensive logging with structured logs
- Graceful shutdown handling
- Docker and Dev Container support
- VSCode debugging integration
- Database migrations

## Requirements

- Go 1.21 or higher
- Docker (optional, for dev containers)
- PostgreSQL (optional, for production setup)

## Quick Start

### Local Development

1. Clone the repository:
```bash
git clone <repository-url>
cd go-backend
```

2. Copy the example environment file:
```bash
cp .env.example .env
```

3. Run the application:
```bash
make run
```

The API will be available at `http://localhost:8080`

### Using Dev Containers

1. Open the project in VSCode
2. When prompted, click "Reopen in Container"
3. Once the container is built, run:
```bash
make run-dev
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check endpoint |
| GET | `/api/v1/projects` | List all projects |
| GET | `/api/v1/projects/{id}` | Get a single project |
| POST | `/api/v1/projects` | Create a new project |
| PUT | `/api/v1/projects/{id}` | Update an existing project |
| DELETE | `/api/v1/projects/{id}` | Delete a project |

See [docs/API.md](docs/API.md) for detailed API documentation.

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
- `DATABASE_TYPE`: `sqlite` or `postgres`
- `DATABASE_URL`: Database connection string
- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: `debug`, `info`, `warn`, or `error`

## Development

### Running Tests

```bash
# Run unit tests
make test

# Run HTTP API tests (requires VSCode REST Client extension)
# First start the server: make run
# Then open files in tests/http/ and click "Send Request"
```

### Running with Hot Reload

```bash
make run-dev
```

### Linting

```bash
make lint
```

### Building

```bash
make build
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

## Debugging

The project includes VSCode debug configurations:

1. **Launch API**: Debug with SQLite
2. **Launch API (PostgreSQL)**: Debug with PostgreSQL
3. **Debug Current Test**: Debug the current test file

To debug:
1. Set breakpoints in your code
2. Press F5 or go to Run > Start Debugging
3. Select the appropriate configuration

## Production Deployment

For production deployment:

1. Set `DATABASE_TYPE=postgres`
2. Configure `DATABASE_URL` with your PostgreSQL connection string
3. Set `ENVIRONMENT=production`
4. Set `LOG_LEVEL=info` or `warn`

## Contributing

This is an educational project demonstrating Go best practices. Key principles:

- Keep it simple and readable
- Follow standard Go patterns
- Document complex logic
- Write tests for critical paths
- Use the standard library where possible

## License

This project is for educational purposes.