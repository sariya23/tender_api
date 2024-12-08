-- +goose Up
-- +goose StatementBegin
alter table organization 
rename column organization_type to organization_type_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
