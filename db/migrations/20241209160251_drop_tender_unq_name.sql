-- +goose Up
-- +goose StatementBegin
alter table tender drop constraint tender_name_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
