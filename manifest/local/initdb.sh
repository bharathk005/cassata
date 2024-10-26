#!/bin/bash

set -e

until pg_isready -d postgres -U postgres; do
  echo "Waiting for database to be ready..."
  sleep 2
done

# Connect to the database and create tables
psql -U postgres -d postgres << EOF

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create workspaces table
CREATE TABLE IF NOT EXISTS workspaces (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    resource VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create user_workspace junction table
CREATE TABLE IF NOT EXISTS user_workspace (
    user_id INTEGER REFERENCES users(id),
    workspace_id INTEGER REFERENCES workspaces(id),
    PRIMARY KEY (user_id, workspace_id)
);

-- Create workspace_permissions junction table
CREATE TABLE IF NOT EXISTS workspace_permissions (
    workspace_id INTEGER REFERENCES workspaces(id),
    permission_id INTEGER REFERENCES permissions(id),
    PRIMARY KEY (workspace_id, permission_id)
);

CREATE TABLE IF NOT EXISTS resource_maps (
    id SERIAL PRIMARY KEY,
    provider VARCHAR(255) NOT NULL,
    resource_group VARCHAR(255) NOT NULL,
    resource_type VARCHAR(255) NOT NULL,
    k8s_api_group VARCHAR(255) NOT NULL,
    k8s_api_version VARCHAR(255) NOT NULL,
    k8s_resource VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);


-- Add unique constraints to prevent duplicate entries

-- Add unique constraint to users table
ALTER TABLE users ADD CONSTRAINT unique_user_name UNIQUE (name);

-- Add unique constraint to workspaces table
ALTER TABLE workspaces ADD CONSTRAINT unique_workspace_name UNIQUE (name);

-- Add unique constraint to permissions table
ALTER TABLE permissions ADD CONSTRAINT unique_permission UNIQUE (resource, action);

-- Add unique constraint to user_workspace table
ALTER TABLE user_workspace ADD CONSTRAINT unique_user_workspace UNIQUE (user_id, workspace_id);

-- Add unique constraint to workspace_permissions table
ALTER TABLE workspace_permissions ADD CONSTRAINT unique_workspace_permission UNIQUE (workspace_id, permission_id);

-- Add unique constraint to resource_maps table
ALTER TABLE resource_maps ADD CONSTRAINT unique_resource_map UNIQUE (provider, resource_group, resource_type);

-- Insert admin user
INSERT INTO users (name, password_hash) 
VALUES ('admin', '\$2y\$10\$Nuq.6WVIgxY9KuR2nl2Id.TvYRc4tm4iQCd1RKiQUrp5loCcUNHU2') 
ON CONFLICT (name) DO NOTHING;

-- Insert admin workspace
INSERT INTO workspaces (name) 
VALUES ('admin') 
ON CONFLICT (name) DO NOTHING;

-- Insert admin permission
INSERT INTO permissions (resource, action) 
VALUES ('/admin/users', 'Create') 
ON CONFLICT (resource, action) DO NOTHING;

-- Insert admin permission
INSERT INTO permissions (resource, action) 
VALUES ('/admin/users', 'Create'), ('/admin/workspaces', 'Create'), ('/admin/permissions', 'Create'), ('/admin/permission-assignments', 'Create')
ON CONFLICT (resource, action) DO NOTHING;

-- Link admin user to admin workspace
INSERT INTO user_workspace (user_id, workspace_id)
SELECT u.id, w.id 
FROM users u, workspaces w 
WHERE u.name = 'admin' AND w.name = 'admin'
ON CONFLICT DO NOTHING;

-- Link admin workspace to admin permission
INSERT INTO workspace_permissions (workspace_id, permission_id)
SELECT w.id, p.id 
FROM workspaces w, permissions p 
WHERE w.name = 'admin' AND p.resource = '/admin/users' AND p.action = 'Create'
ON CONFLICT DO NOTHING;

INSERT INTO workspace_permissions (workspace_id, permission_id)
SELECT w.id, p.id 
FROM workspaces w, permissions p 
WHERE w.name = 'admin' AND p.resource = '/admin/workspaces' AND p.action = 'Create'
ON CONFLICT DO NOTHING;

INSERT INTO workspace_permissions (workspace_id, permission_id)
SELECT w.id, p.id 
FROM workspaces w, permissions p 
WHERE w.name = 'admin' AND p.resource = '/admin/permissions' AND p.action = 'Create'
ON CONFLICT DO NOTHING;

INSERT INTO workspace_permissions (workspace_id, permission_id)
SELECT w.id, p.id 
FROM workspaces w, permissions p 
WHERE w.name = 'admin' AND p.resource = '/admin/permission-assignments' AND p.action = 'Create'
ON CONFLICT DO NOTHING;

-- Insert GCP resource map
INSERT INTO resource_maps (provider, resource_group, resource_type, k8s_api_group, k8s_api_version, k8s_resource) 
VALUES ('gcp', 'iam', 'ServiceAccount', 'iam.gcp.crossplane.io', 'v1alpha1', 'serviceaccounts') 
ON CONFLICT (provider, resource_group, resource_type) DO NOTHING;

EOF

echo "Database and tables created successfully."

