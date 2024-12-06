-- +goose Up
-- +goose StatementBegin
create table if not exists employee (
    employee_id bigint generated always as identity primary key,
    username varchar(50) unique not null,
    first_name varchar(50),
    last_name varchar(50),
    created_at timestamp default CURRENT_TIMESTAMP,
    updated_at timestamp default CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists employee;
-- +goose StatementEnd
