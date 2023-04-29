-- +goose Up
-- +goose StatementBegin
CREATE INDEX "code_index"
ON "currency" ("code");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX "codeIndex";
-- +goose StatementEnd
