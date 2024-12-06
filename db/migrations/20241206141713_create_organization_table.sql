-- +goose Up
-- +goose StatementBegin
create table if not exists organization (
    organization_id bigint generated always as identity primary key,
    name varchar(100) not null,
    description text,
    type organization_type,
    created_at timestamp default CURRENT_TIMESTAMP,
    updated_at timestamp default CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists organization;
-- +goose StatementEnd
