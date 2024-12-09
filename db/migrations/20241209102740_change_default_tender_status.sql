-- +goose Up
-- +goose StatementBegin
alter table tender
alter column status
set default 'CREATED';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
