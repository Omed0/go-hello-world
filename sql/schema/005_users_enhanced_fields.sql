-- +goose Up
ALTER TABLE users
ADD COLUMN password_hash VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN age INTEGER CHECK (age >= 0 AND age <= 150),
ADD COLUMN gender VARCHAR(20) CHECK (gender IN ('male', 'female', 'other', 'prefer_not_to_say')),
ADD COLUMN role VARCHAR(50) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin', 'moderator')),
ADD COLUMN organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE users 
DROP COLUMN password_hash,
DROP COLUMN age,
DROP COLUMN gender,
DROP COLUMN role,
DROP COLUMN organization_id;
