-- +goose Up
-- +goose StatementBegin
CREATE TABLE oauth_tokens (
    id SERIAL PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL,
    code TEXT NOT NULL,
    "access" TEXT NOT NULL,
    refresh TEXT NOT NULL,
    data JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_oauth_tokens_expires_at ON oauth_tokens (expires_at);
CREATE INDEX IF NOT EXISTS idx_oauth_tokens_code ON oauth_tokens (code);
CREATE INDEX IF NOT EXISTS idx_oauth_tokens_access ON oauth_tokens (access);
CREATE INDEX IF NOT EXISTS idx_oauth_tokens_refresh ON oauth_tokens (refresh);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS oauth_tokens;
-- +goose StatementEnd
