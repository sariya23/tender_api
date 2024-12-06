-- +goose Up
-- +goose StatementBegin
create table if not exists nsi_organization_type (
    nsi_organization_type_id smallint primary key,
    type varchar(50) not null unique
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists nsi_organization_type;
-- +goose StatementEnd
