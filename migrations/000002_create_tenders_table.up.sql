CREATE TABLE tenders (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    service_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    organization_id UUID NOT NULL,
    version INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);