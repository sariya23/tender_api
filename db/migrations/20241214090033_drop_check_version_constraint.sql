-- +goose Up
-- +goose StatementBegin
alter table tender
drop constraint selected_version_check;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
