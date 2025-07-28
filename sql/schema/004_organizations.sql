-- +goose Up
CREATE TABLE organizations (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_organizations_name ON organizations(name);
CREATE INDEX idx_organizations_deleted_at ON organizations(deleted_at) WHERE deleted_at IS NULL;

-- Insert default organization
INSERT INTO organizations (id, name, description) 
VALUES (
    '00000000-0000-0000-0000-000000000000',
    'Default Organization',
    'Default organization for users without a specific organization'
);

-- +goose Down
DROP TABLE organizations;
