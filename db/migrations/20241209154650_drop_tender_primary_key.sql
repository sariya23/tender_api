-- +goose Up
-- +goose StatementBegin
alter table tender 
drop column tender_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
