-- +goose Up
-- +goose StatementBegin
insert into nsi_organization_type values
(1, 'IE'),
(2, 'LLC'),
(3, 'JSC');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from nsi_organization_type where nsi_organization_type_id in (1, 2, 3);
-- +goose StatementEnd
