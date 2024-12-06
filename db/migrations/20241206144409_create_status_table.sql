-- +goose Up
-- +goose StatementBegin
create table if not exists nsi_tender_status (
    nsi_tender_status_id smallint primary key,
    status varchar(50) not null unique
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists nsi_tender_status;
-- +goose StatementEnd
