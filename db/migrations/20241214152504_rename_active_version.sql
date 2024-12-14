-- +goose Up
-- +goose StatementBegin
alter table tender
rename column selected_version to is_active_version;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
