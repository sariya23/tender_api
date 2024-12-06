-- +goose Up
-- +goose StatementBegin
insert into nsi_tender_status values
(1, 'CREATED'),
(2, 'PUBLISHED'),
(3, 'CLOSED');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from nsi_tender_status where nsi_tender_status_id in (1, 2, 3);
-- +goose StatementEnd
