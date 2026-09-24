-- name: ListSystemSettings :many
SELECT *
FROM system_settings
ORDER BY category ASC, key ASC;

-- name: ListSystemSettingsByCategory :many
SELECT *
FROM system_settings
WHERE category = $1
ORDER BY key ASC;

-- name: GetSystemSetting :one
SELECT *
FROM system_settings
WHERE category = $1 AND key = $2;

-- name: UpsertSystemSetting :one
INSERT INTO system_settings (
  category,
  key,
  value,
  is_secret,
  updated_by
) VALUES (
  $1, $2, $3, $4, $5
)
ON CONFLICT (category, key) DO UPDATE
  SET value = EXCLUDED.value,
      is_secret = EXCLUDED.is_secret,
      updated_by = EXCLUDED.updated_by
RETURNING *;

-- name: DeleteSystemSetting :exec
DELETE FROM system_settings
WHERE category = $1 AND key = $2;

-- name: DeleteSystemSettingsByCategory :exec
DELETE FROM system_settings
WHERE category = $1;
