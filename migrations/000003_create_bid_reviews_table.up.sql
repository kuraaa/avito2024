CREATE TABLE bid_reviews (
    id UUID PRIMARY KEY,
    description TEXT NOT NULL,
    tender_id UUID NOT NULL REFERENCES tenders(id),
    author_username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL
);