-- +goose Up
-- +goose StatementBegin
CREATE TABLE specialities(
    id SERIAL PRIMARY KEY,
    speciality VARCHAR(255) NOT NULL
);

INSERT INTO specialities (speciality) VALUES
('Маляр'),
('Штукатурщик');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE specialities;
-- +goose StatementEnd
