-- +goose Up
-- +goose StatementBegin
CREATE TABLE workers (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    fullname VARCHAR(255),
    phone VARCHAR(25),
    speciality VARCHAR(255),
    city VARCHAR(60)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workers;
-- +goose StatementEnd
