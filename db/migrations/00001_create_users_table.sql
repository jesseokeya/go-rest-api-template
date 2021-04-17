-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_role AS ENUM ('', 'member', 'admin', 'blocked');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    password_hash varchar(60) NOT NULL DEFAULT '',
    email text NOT NULL DEFAULT '',
    role user_role NOT NULL DEFAULT '',
    first_name text NOT NULL DEFAULT '',
    last_name text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
