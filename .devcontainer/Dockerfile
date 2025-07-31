# Start from the official Go development container image
# This provides a consistent development environment
FROM mcr.microsoft.com/devcontainers/go:1-1.24-bullseye

# Install additional development tools
# These tools enhance the development experience
RUN apt-get update && apt-get install -y \
    postgresql-client \
    make \
    && rm -rf /var/lib/apt/lists/*

# Install Go development tools in a temporary module context
# air: hot reload for Go applications
# golangci-lint: comprehensive Go linter
# delve: Go debugger
RUN cd /tmp && \
    go mod init temp && \
    go install github.com/air-verse/air@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    rm -rf /tmp/go.mod /tmp/go.sum

# Ensure vscode user has sudo access
# The base image already includes the vscode user
RUN usermod -aG sudo vscode && \
    echo "vscode ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# Set the working directory
WORKDIR /workspace

# Switch to the non-root user
USER vscode