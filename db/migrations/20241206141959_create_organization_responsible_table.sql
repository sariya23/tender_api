-- +goose Up
-- +goose StatementBegin
create table organization_responsible (
    organization_responsible_id bigint generated always as identity primary key,
    organization_id bigint references organization(organization_id) on delete cascade,
    employee_id bigint references employee(employee_id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists organization_responsible cascade;
-- +goose StatementEnd
