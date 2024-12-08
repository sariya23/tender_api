-- +goose Up
-- +goose StatementBegin
alter table tender 
rename column status to status_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
