-- +goose Up
-- +goose StatementBegin
drop table tender;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
