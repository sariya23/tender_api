-- +goose Up
-- +goose StatementBegin
create unique index unique_selected_version_per_tender
on tender (tender_id)
where selected_version = true;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
