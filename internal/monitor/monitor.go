package monitor

import (
	"context"
	"fmt"
	"time"

	storage "github.com/grahovsky/system-stats-daemon/internal/storage/memory"
)

func New(ctx context.Context) {
	sm := storage.New(20)

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

	go func() {
		for {
			select {
			case <-tiker.C:
				sm.Push(d, time.Now())
				num++
				d.cpu -= 1
				d.mem++
				if num >= 5 {
					fmt.Println(num)
					sm.Print()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
