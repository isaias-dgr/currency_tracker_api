-- +goose Up
-- +goose StatementBegin
ALTER TABLE currency
ADD COLUMN code VARCHAR(5),
ADD COLUMN value FLOAT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE currency
DROP COLUMN code,
DROP COLUMN value;
-- +goose StatementEnd