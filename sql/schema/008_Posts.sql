-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,
    title TEXT NOT NULL ,
    url VARCHAR(512) UNIQUE NOT NULL,
    description TEXT ,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE NOT NULL
);

-- +goose Down
DROP TABLE posts;