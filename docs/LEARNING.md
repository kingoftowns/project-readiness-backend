# Learning Resources

This document provides resources and explanations for learning Go through this project.

## Go Concepts Demonstrated

### 1. Project Structure

This project follows the standard Go project layout:

```
cmd/        # Main applications
internal/   # Private application code
migrations/ # Database migrations
docs/       # Documentation
```

**Learn more**: [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

### 2. Interfaces

Interfaces define contracts that types must fulfill:

```go
type GitLabRepository interface {
    Create(ctx context.Context, project *models.Project) error
    GetByID(ctx context.Context, projectID string) (*models.Project, error)
}
```

**Key concepts**:
- Implicit interface satisfaction
- Interface segregation
- Dependency inversion

**Learn more**: [Effective Go - Interfaces](https://go.dev/doc/effective_go#interfaces)

### 3. Error Handling

Go's explicit error handling:

```go
project, err := h.repo.GetByID(ctx, projectID)
if err != nil {
    if err.Error() == "project not found" {
        h.respondWithError(w, http.StatusNotFound, "Project not found")
        return
    }
    // Handle other errors
}
```

**Key concepts**:
- Errors as values
- Error wrapping with `fmt.Errorf`
- Custom error types

**Learn more**: [Error handling in Go](https://go.dev/blog/error-handling-and-go)

### 4. Context

Context for cancellation and request-scoped values:

```go
func (r *gitLabRepo) Create(ctx context.Context, project *models.Project) error
```

**Key concepts**:
- Request cancellation
- Timeouts
- Request-scoped values

**Learn more**: [Go Concurrency Patterns: Context](https://go.dev/blog/context)

### 5. Struct Tags

Struct tags for JSON and database mapping:

```go
type Project struct {
    ProjectID string `json:"project_id" db:"project_id"`
}
```

**Learn more**: [The Go Programming Language Specification - Struct types](https://go.dev/ref/spec#Struct_types)

### 6. Embedding

Embedding types for composition:

```go
type responseWriter struct {
    http.ResponseWriter  // Embedded interface
    statusCode int
}
```

**Learn more**: [Effective Go - Embedding](https://go.dev/doc/effective_go#embedding)

## Best Practices Demonstrated

### 1. Dependency Injection

```go
func NewGitLabHandler(repo repository.GitLabRepository, logger *slog.Logger) *GitLabHandler {
    return &GitLabHandler{
        repo:   repo,
        logger: logger,
    }
}
```

**Benefits**:
- Testability
- Loose coupling
- Explicit dependencies

### 2. Configuration Management

```go
func Load() (*Config, error) {
    cfg := &Config{
        Port: getEnv("PORT", "8080"),
        // ...
    }
}
```

**Benefits**:
- 12-factor app compliance
- Environment-specific settings
- Default values

### 3. Graceful Shutdown

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
srv.Shutdown(ctx)
```

**Benefits**:
- Clean resource cleanup
- In-flight request completion
- Data integrity

## Testing Patterns

### 1. Table-Driven Tests

```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {"valid input", "test", "TEST", false},
    {"empty input", "", "", true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test logic
    })
}
```

### 2. Mock Interfaces

```go
type mockRepository struct {
    projects map[string]*models.Project
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*models.Project, error) {
    project, ok := m.projects[id]
    if !ok {
        return nil, errors.New("not found")
    }
    return project, nil
}
```

## Recommended Learning Path

1. **Go Tour**: Start with [A Tour of Go](https://go.dev/tour/)
2. **Effective Go**: Read [Effective Go](https://go.dev/doc/effective_go)
3. **Go by Example**: Work through [Go by Example](https://gobyexample.com/)
4. **This Project**: Study this codebase in order:
   - Models (domain logic)
   - Repository (data access)
   - Handlers (HTTP layer)
   - Main (application bootstrap)

## Useful Commands

### Development
```bash
# Run tests
go test ./...

# Run with race detector
go test -race ./...

# Check for issues
go vet ./...

# Format code
go fmt ./...
```

### Debugging
```bash
# Run with debug output
GODEBUG=gctrace=1 go run cmd/api/main.go

# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=.

# Analyze profile
go tool pprof cpu.prof
```

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Go Blog](https://go.dev/blog/)
- [Go Wiki](https://github.com/golang/go/wiki)
- [Awesome Go](https://github.com/avelino/awesome-go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## Community

- [Go Forum](https://forum.golangbridge.org/)
- [Gophers Slack](https://invite.slack.golangbridge.org/)
- [r/golang](https://www.reddit.com/r/golang/)

## Books

1. "The Go Programming Language" by Donovan and Kernighan
2. "Go in Action" by William Kennedy
3. "Concurrency in Go" by Katherine Cox-Buday
4. "Learning Go" by Jon Bodner