-- +goose Up
-- +goose StatementBegin
alter table organization 
add organization_type smallint not null references nsi_tender_status(nsi_tender_status_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table organization
drop column organization_type;
-- +goose StatementEnd
