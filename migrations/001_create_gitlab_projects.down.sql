-- Drop the gitlab_projects table and its associated index
DROP INDEX IF EXISTS idx_gitlab_projects_created_at;
DROP TABLE IF EXISTS gitlab_projects;