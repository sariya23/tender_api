-- +goose Up
-- +goose StatementBegin
insert into nsi_tender_status values 
(1, 'CREATED'),
(2, 'PUBLISHED'),
(3, 'CLOSED');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
