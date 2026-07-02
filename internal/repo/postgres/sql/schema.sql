CREATE TYPE issue_status AS ENUM (
    'NEW',
    'ASSIGNED',
    'IN_PROGRESS',
    'RESOLVED',
    'CLOSED',
    'REOPENED'
);
 
CREATE TYPE issue_resolution AS ENUM (
    'FIXED',
    'INVALID',
    'WONTFIX',
    'WORKSFORME',
    'Unspecified'
);
 
CREATE TYPE issue_type AS ENUM (
    'COSMETIC',
    'BUG',
    'FEATURE',
    'PERFORMANCE'
);
 
CREATE TYPE issue_priority AS ENUM (
    'CRITICAL',
    'MAJOR',
    'IMPORTANT',
    'MINOR'
);
 
CREATE TABLE IF NOT EXISTS users (
    id            TEXT PRIMARY KEY,
    first_name    TEXT NOT NULL,
    last_name     TEXT NOT NULL,
    email_address TEXT NOT NULL UNIQUE
);
 
CREATE TABLE IF NOT EXISTS projects (
    id            TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT ''
);
 
CREATE TABLE IF NOT EXISTS issues (
    id            TEXT PRIMARY KEY,
    create_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modify_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 
    summary     TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
 
    status      issue_status     NOT NULL DEFAULT 'NEW',
    resolution  issue_resolution NOT NULL DEFAULT 'Unspecified',
    type        issue_type       NOT NULL DEFAULT 'BUG',
    priority    issue_priority         NOT NULL DEFAULT 'MINOR',
 
    project_id  TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    assignee_id TEXT          REFERENCES users(id)    ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_issues_project_id  ON issues(project_id);
CREATE INDEX IF NOT EXISTS idx_issues_assignee_id ON issues(assignee_id);
CREATE INDEX IF NOT EXISTS idx_issues_status      ON issues(status);
CREATE INDEX IF NOT EXISTS idx_issues_priority    ON issues(priority);
 
