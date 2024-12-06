-- +goose Up
-- +goose StatementBegin
drop type if exists organization_type;
CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop type if exists organization_type;
-- +goose StatementEnd
