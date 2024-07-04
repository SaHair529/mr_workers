-- +goose Up
-- +goose StatementBegin
ALTER TABLE requests DROP CONSTRAINT IF EXISTS requests_telegram_id_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE requests ADD CONSTRAINT requests_telegram_id_key UNIQUE (telegram_id);
-- +goose StatementEnd
