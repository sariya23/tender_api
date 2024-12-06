-- +goose Up
-- +goose StatementBegin
alter table organization
drop column type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
