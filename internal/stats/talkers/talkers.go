package talkers

import (
	"context"

	"github.com/grahovsky/system-stats-daemon/internal/config"
	"github.com/grahovsky/system-stats-daemon/internal/storage"
)

func GetStats(ctx context.Context, st storage.Storage) {
	if !config.Settings.Stats.Collect.Talkers {
		return
	}
	GetStatsOs(ctx, st)
}
