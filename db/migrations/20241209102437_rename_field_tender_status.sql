-- +goose Up
-- +goose StatementBegin
alter table tender
rename column status_id to status;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
