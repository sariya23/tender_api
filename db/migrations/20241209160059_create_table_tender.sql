-- +goose Up
-- +goose StatementBegin
create table if not exists tender (
    tender_id bigint generated always as identity,
    name text not null unique,
    description text,
    service_type text not null,
    status varchar(9) not null check (status in ('CREATED', 'CLOSED', 'PUBLISHED')) default 'CREATED',
    organization_id bigint not null references organization(organization_id),
    creator_username text not null,
    version int not null default 1 check(version > 0),
    primary key (name, version)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists tender;
-- +goose StatementEnd