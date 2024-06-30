-- +goose Up
-- +goose StatementBegin
CREATE TABLE requests (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    specialist VARCHAR(255),
    city VARCHAR(65),
    description TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION  update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER update_updated_at
BEFORE UPDATE ON requests
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_updated_at ON requests;
-- +goose StatementEnd

-- +goose StatementBegin
DROP FUNCTION IF EXISTS update_updated_at_column();
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE requests;
-- +goose StatementEnd
