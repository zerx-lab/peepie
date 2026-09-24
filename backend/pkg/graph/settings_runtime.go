package graph

import (
	"context"

	"pentagi/pkg/config"
	"pentagi/pkg/database"
)

// applySetting persists one (category, key) -> value pair to system_settings
// and immediately applies it to cfg.Overrides, so every subsequent read
func applySetting(ctx context.Context, db database.Querier, cfg *config.Config, uid int64, category, key, value string, isSecret bool) error {
	_, err := db.UpsertSystemSetting(ctx, database.UpsertSystemSettingParams{
		Category:  category,
		Key:       key,
		Value:     value,
		IsSecret:  isSecret,
		UpdatedBy: database.Int64ToNullInt64(&uid),
	})
	if err != nil {
		return err
	}

	cfg.Overrides.Set(category, key, value)
	return nil
}

// applySecretSetting is applySetting for a *string input that must not clear
// an already-configured secret when the caller omits it (nil): the frontend
// never re-sends a plaintext secret it only echoed as "configured", so a nil
// input here means "leave unchanged", not "clear".
func applySecretSetting(ctx context.Context, db database.Querier, cfg *config.Config, uid int64, category, key string, value *string) error {
	if value == nil {
		return nil
	}
	return applySetting(ctx, db, cfg, uid, category, key, *value, true)
}
