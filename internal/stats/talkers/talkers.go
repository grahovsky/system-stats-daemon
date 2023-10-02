package talkers

import (
	"context"

	"github.com/grahovsky/system-stats-daemon/internal/storage"
)

func GetStats(ctx context.Context, st storage.Storage) {
	GetStatsOs(ctx, st)
}
