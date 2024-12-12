-- +goose Up
-- +goose StatementBegin
alter table tender
add constraint selected_version_check unique(tender_id, selected_version);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
