package monitor

import (
	"context"
	"time"

	"github.com/grahovsky/system-stats-daemon/internal/storage"
	memoryStorage "github.com/grahovsky/system-stats-daemon/internal/storage/memory"
)

func New(ctx context.Context) {
	var sm storage.Storage = memoryStorage.New()
	d := struct {
		cpu int64
		mem int64
	}{
		cpu: 100,
		mem: 16,
	}

	num := 0
	tiker := time.NewTicker(1 * time.Second)
	defer tiker.Stop()
	defer sm.Show()

	for {
		select {
		case <-tiker.C:
			sm.Push(d, time.Now())
			num++
			d.cpu -= 1
			d.mem++
			if num >= 5 {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
