-- +goose Up
-- +goose StatementBegin
alter table tender
add column selected_version bool default true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
