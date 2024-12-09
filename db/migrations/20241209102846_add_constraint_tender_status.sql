-- +goose Up
-- +goose StatementBegin
alter table tender
add constraint check_tender_status
check (status in ('CREATED', 'CLOSED', 'PUBLISHED'))
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
