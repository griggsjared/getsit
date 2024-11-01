-- +goose Up
-- +goose StatementBegin
CREATE TABLE url_entries (
    id SERIAL PRIMARY KEY,
    url TEXT UNIQUE DEFAULT NULL,
    token TEXT UNIQUE DEFAULT NULL,
    visit_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE url_entries;
-- +goose StatementEnd
