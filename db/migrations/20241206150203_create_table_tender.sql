-- +goose Up
-- +goose StatementBegin
create table if not exists tender (
    tender_id bigint generated always as identity primary key,
    name text not null unique,
    description text,
    service_type text not null,
    status smallint not null references nsi_tender_status(nsi_tender_status_id) default 1,
    organization_id bigint not null references organization(organization_id),
    creator_username text not null,
    version int not null default 1 check(version > 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists tender;
-- +goose StatementEnd
