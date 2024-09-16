CREATE TABLE bids (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    tender_id UUID NOT NULL REFERENCES tenders(id),
    author_type VARCHAR(50) NOT NULL,
    author_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    version INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);