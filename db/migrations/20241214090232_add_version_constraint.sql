-- +goose Up
-- +goose StatementBegin
alter table tender add constraint unique_tender_version UNIQUE (tender_id, version);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
