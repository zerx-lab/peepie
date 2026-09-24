-- +goose Up
-- +goose StatementBegin
CREATE TABLE system_settings (
  id          BIGINT       PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  category    TEXT         NOT NULL,
  key         TEXT         NOT NULL,
  value       TEXT         NOT NULL,
  is_secret   BOOLEAN      NOT NULL DEFAULT FALSE,
  updated_by  BIGINT       NULL REFERENCES users(id) ON DELETE SET NULL,
  created_at  TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,
  updated_at  TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT system_settings_category_key_unique UNIQUE (category, key)
);

CREATE INDEX system_settings_category_idx ON system_settings(category);

CREATE TRIGGER update_system_settings_modified
  BEFORE UPDATE ON system_settings
  FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- Add privileges for Admin role (role_id = 1) — only admins may view/edit
-- runtime system settings (LLM provider keys, search engine keys, etc. are
-- as sensitive as the .env file they replace).
INSERT INTO privileges (role_id, name) VALUES
    (1, 'settings.system.admin'),
    (1, 'settings.system.view'),
    (1, 'settings.system.edit')
    ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM privileges WHERE name IN (
  'settings.system.admin',
  'settings.system.view',
  'settings.system.edit'
);

DROP TABLE IF EXISTS system_settings;
-- +goose StatementEnd
