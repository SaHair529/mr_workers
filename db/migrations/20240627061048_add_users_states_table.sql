-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_states (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    state VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_states;
-- +goose StatementEnd
