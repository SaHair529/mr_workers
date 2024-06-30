-- +goose Up
-- +goose StatementBegin
ALTER TABLE requests
ADD COLUMN free BOOLEAN DEFAULT TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE requests
DROP COLUMN free;
-- +goose StatementEnd
