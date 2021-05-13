-- +goose Up
-- +goose StatementBegin
CREATE TABLE oauth_clients (
    id text NOT NULL,
    secret text NOT NULL,
    domain text NOT NULL,
    data JSONB NOT NULL DEFAULT '{}'
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS oauth_clients;
-- +goose StatementEnd
