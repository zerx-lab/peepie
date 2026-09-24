package database

import "context"

// LoadSystemSettingsOverrides reads every row from system_settings and groups
// it into category -> key -> value, the shape config.Overrides.Load expects.
// Kept here (rather than in pkg/config) because pkg/config must not import
// pkg/database.
func LoadSystemSettingsOverrides(ctx context.Context, q Querier) (map[string]map[string]string, error) {
	rows, err := q.ListSystemSettings(ctx)
	if err != nil {
		return nil, err
	}

	out := make(map[string]map[string]string)
	for _, row := range rows {
		if out[row.Category] == nil {
			out[row.Category] = make(map[string]string)
		}
		out[row.Category][row.Key] = row.Value
	}

	return out, nil
}
