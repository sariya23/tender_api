-- +goose Up
-- +goose StatementBegin
alter table organization 
add organization_type smallint not null references nsi_organization_type(nsi_organization_type_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table organization
drop column organization_type;
-- +goose StatementEnd
