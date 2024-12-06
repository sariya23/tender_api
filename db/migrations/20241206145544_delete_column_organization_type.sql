-- +goose Up
-- +goose StatementBegin
alter table organization
drop column organization_type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
