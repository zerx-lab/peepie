-- +goose Up
-- +goose StatementBegin
-- Rename the seeded administrator login to the peepie product address.
-- Password and password_change_required are left untouched.
UPDATE users
    SET mail = 'admin@peepie.com'
    WHERE mail = 'admin@pentagi.com'
      AND NOT EXISTS (SELECT 1 FROM users WHERE mail = 'admin@peepie.com');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE users
    SET mail = 'admin@pentagi.com'
    WHERE mail = 'admin@peepie.com'
      AND NOT EXISTS (SELECT 1 FROM users WHERE mail = 'admin@pentagi.com');
-- +goose StatementEnd
