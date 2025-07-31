# Architecture Documentation

## Overview

This application follows Clean Architecture principles, emphasizing separation of concerns, dependency inversion, and testability.

## Architecture Layers

### 1. Domain Layer (`internal/models`)

The core business logic and entities. This layer has no dependencies on external packages or other layers.

- **Project**: The main domain entity representing a GitLab project and its readiness checks
- **Business Rules**: Methods like `IsProductionReady()` and `FailedChecks()` encapsulate business logic

### 2. Repository Layer (`internal/repository`)

Handles data persistence and retrieval. This layer implements the repository pattern to abstract database operations.

- **Interface**: `GitLabRepository` defines the contract for data operations
- **Implementation**: `gitLabRepo` provides the concrete database implementation
- **Benefits**: Easy to mock for testing, database-agnostic interface

### 3. Handler Layer (`internal/handlers`)

Contains HTTP request handlers that orchestrate between the HTTP layer and the business logic.

- **GitLabHandler**: Handles all GitLab project-related HTTP endpoints
- **Responsibilities**: Request validation, response formatting, error handling
- **Dependencies**: Injected repository interface for data access

### 4. Infrastructure Layer

#### Configuration (`internal/config`)
- Environment-based configuration following 12-factor app principles
- Validation of configuration values
- Default values for development

#### Database (`internal/database`)
- Database connection management
- Support for multiple database types (SQLite, PostgreSQL)
- Migration system for schema management

#### Router (`internal/router`)
- HTTP routing using chi router
- Middleware configuration
- Request logging and monitoring

## Design Patterns

### Repository Pattern
Abstracts data access logic, making the application database-agnostic:

```go
type GitLabRepository interface {
    Create(ctx context.Context, project *models.Project) error
    GetByID(ctx context.Context, projectID string) (*models.Project, error)
    // ...
}
```

### Dependency Injection
Dependencies are injected through constructors, improving testability:

```go
func NewGitLabHandler(repo repository.GitLabRepository, logger *slog.Logger) *GitLabHandler
```

### Interface Segregation
Small, focused interfaces that are easy to implement and mock:

```go
// Each handler method handles one specific HTTP endpoint
func (h *GitLabHandler) GetProject(w http.ResponseWriter, r *http.Request)
```

## Data Flow

1. **Request**: HTTP request arrives at the router
2. **Middleware**: Request passes through middleware (logging, timeout, etc.)
3. **Handler**: Router directs request to appropriate handler
4. **Repository**: Handler calls repository for data operations
5. **Database**: Repository executes SQL queries
6. **Response**: Handler formats and returns response

## Error Handling

Errors are handled at each layer with appropriate context:

1. **Database errors**: Wrapped with context in repository layer
2. **Business errors**: Returned as specific error types
3. **HTTP errors**: Converted to appropriate HTTP status codes

## Configuration Management

Following 12-factor app principles:

- Configuration via environment variables
- No hardcoded values
- Sensible defaults for development
- Validation at startup

## Database Design

### Multi-Database Support
- SQLite for development and testing
- PostgreSQL for production
- Same codebase works with both

### Migration System
- SQL-based migrations
- Up and down migrations
- Version tracking in database

## Security Considerations

1. **Input Validation**: All user input is validated
2. **SQL Injection Protection**: Using parameterized queries
3. **Error Messages**: Generic error messages to avoid information leakage
4. **Timeouts**: Request timeouts to prevent DoS

## Scalability

The architecture supports horizontal scaling:

1. **Stateless**: No session state in the application
2. **Database Connection Pooling**: Efficient connection management
3. **Context-based Cancellation**: Proper request lifecycle management

## Testing Strategy

### Unit Tests
- Repository layer with in-memory SQLite
- Handlers with mock repositories
- Models with pure business logic tests

### Integration Tests
- Full API tests with test database
- Migration tests

### Example Test Structure
```go
func TestGitLabRepository_Create(t *testing.T) {
    // Setup test database
    // Create repository
    // Test create operation
    // Assert results
}
```

## Deployment Architecture

### Development
- Single binary with SQLite
- Hot reload with air
- Debug mode with verbose logging

### Production
- Docker container deployment
- PostgreSQL database
- Environment-based configuration
- Structured JSON logging

## Future Enhancements

1. **API Versioning**: Support for multiple API versions
2. **Authentication**: JWT or OAuth2 integration
3. **Caching**: Redis for frequently accessed data
4. **Message Queue**: Async processing of checks
5. **Metrics**: Prometheus metrics for monitoring