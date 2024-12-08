-- +goose Up
-- +goose StatementBegin
truncate nsi_tender_status cascade;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
