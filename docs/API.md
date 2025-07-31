# API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently, the API does not require authentication. In a production environment, you would add authentication middleware.

## Endpoints

### Health Check

Check if the API is running and healthy.

**Endpoint:** `GET /api/v1/health`

**Response:**
```json
{
  "status": "healthy",
  "service": "gitlab-readiness-api"
}
```

### List Projects

Get a paginated list of all GitLab projects.

**Endpoint:** `GET /api/v1/projects`

**Query Parameters:**
- `limit` (optional): Number of items per page (default: 50, max: 100)
- `offset` (optional): Number of items to skip (default: 0)

**Response:**
```json
{
  "data": [
    {
      "project_id": "12345",
      "project_present": true,
      "app_name_set": true,
      "moab_id_set": true,
      "codeowners_exists": true,
      "branch_protection_enabled": true,
      "codeowner_approval_required": true,
      "push_merge_restricted": true,
      "force_push_disabled": true,
      "push_rules_enabled": true,
      "min_approvals_required": true,
      "author_approval_prevented": true,
      "committer_approval_prevented": true,
      "approvals_removed_on_commit": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "limit": 50,
    "offset": 0,
    "total": 100
  }
}
```

### Get Project

Get a single project by ID.

**Endpoint:** `GET /api/v1/projects/{id}`

**Response:**
```json
{
  "data": {
    "project": {
      "project_id": "12345",
      "project_present": true,
      "app_name_set": true,
      "moab_id_set": false,
      "codeowners_exists": true,
      "branch_protection_enabled": true,
      "codeowner_approval_required": true,
      "push_merge_restricted": true,
      "force_push_disabled": true,
      "push_rules_enabled": true,
      "min_approvals_required": true,
      "author_approval_prevented": true,
      "committer_approval_prevented": true,
      "approvals_removed_on_commit": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "is_production_ready": false,
    "failed_checks": ["moab_id_set"]
  }
}
```

### Create Project

Create a new project with readiness checks.

**Endpoint:** `POST /api/v1/projects`

**Request Body:**
```json
{
  "project_id": "12345",
  "project_present": true,
  "app_name_set": true,
  "moab_id_set": true,
  "codeowners_exists": true,
  "branch_protection_enabled": true,
  "codeowner_approval_required": true,
  "push_merge_restricted": true,
  "force_push_disabled": true,
  "push_rules_enabled": true,
  "min_approvals_required": true,
  "author_approval_prevented": true,
  "committer_approval_prevented": true,
  "approvals_removed_on_commit": true
}
```

**Response:** `201 Created`
```json
{
  "data": {
    "project_id": "12345",
    "project_present": true,
    "app_name_set": true,
    "moab_id_set": true,
    "codeowners_exists": true,
    "branch_protection_enabled": true,
    "codeowner_approval_required": true,
    "push_merge_restricted": true,
    "force_push_disabled": true,
    "push_rules_enabled": true,
    "min_approvals_required": true,
    "author_approval_prevented": true,
    "committer_approval_prevented": true,
    "approvals_removed_on_commit": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Update Project

Update an existing project's readiness checks.

**Endpoint:** `PUT /api/v1/projects/{id}`

**Request Body:**
```json
{
  "project_present": true,
  "app_name_set": true,
  "moab_id_set": true,
  "codeowners_exists": true,
  "branch_protection_enabled": true,
  "codeowner_approval_required": true,
  "push_merge_restricted": true,
  "force_push_disabled": true,
  "push_rules_enabled": true,
  "min_approvals_required": true,
  "author_approval_prevented": true,
  "committer_approval_prevented": true,
  "approvals_removed_on_commit": true
}
```

**Response:** `200 OK`
```json
{
  "data": {
    "project_id": "12345",
    "project_present": true,
    "app_name_set": true,
    "moab_id_set": true,
    "codeowners_exists": true,
    "branch_protection_enabled": true,
    "codeowner_approval_required": true,
    "push_merge_restricted": true,
    "force_push_disabled": true,
    "push_rules_enabled": true,
    "min_approvals_required": true,
    "author_approval_prevented": true,
    "committer_approval_prevented": true,
    "approvals_removed_on_commit": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Delete Project

Delete a project from the system.

**Endpoint:** `DELETE /api/v1/projects/{id}`

**Response:** `204 No Content`

## Error Responses

All error responses follow this format:

```json
{
  "error": {
    "message": "Description of the error"
  }
}
```

### Common Error Codes

- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists
- `500 Internal Server Error`: Server error

## Example Usage

### Create a new project

```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "12345",
    "project_present": true,
    "app_name_set": true,
    "moab_id_set": true,
    "codeowners_exists": true,
    "branch_protection_enabled": true,
    "codeowner_approval_required": true,
    "push_merge_restricted": true,
    "force_push_disabled": true,
    "push_rules_enabled": true,
    "min_approvals_required": true,
    "author_approval_prevented": true,
    "committer_approval_prevented": true,
    "approvals_removed_on_commit": true
  }'
```

### Get all projects

```bash
curl http://localhost:8080/api/v1/projects?limit=10&offset=0
```

### Get a specific project

```bash
curl http://localhost:8080/api/v1/projects/12345
```

### Update a project

```bash
curl -X PUT http://localhost:8080/api/v1/projects/12345 \
  -H "Content-Type: application/json" \
  -d '{
    "project_present": true,
    "app_name_set": true,
    "moab_id_set": false,
    "codeowners_exists": true,
    "branch_protection_enabled": true,
    "codeowner_approval_required": true,
    "push_merge_restricted": true,
    "force_push_disabled": true,
    "push_rules_enabled": true,
    "min_approvals_required": true,
    "author_approval_prevented": true,
    "committer_approval_prevented": true,
    "approvals_removed_on_commit": true
  }'
```

### Delete a project

```bash
curl -X DELETE http://localhost:8080/api/v1/projects/12345
```