-- +goose Up
-- +goose StatementBegin
alter table tender
alter column status
type varchar(9);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
