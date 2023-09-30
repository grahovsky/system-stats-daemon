package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGRPCServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Microsecond)
	defer cancel()

	go func() {
		srv := NewStatsMonitoringSever(ctx)
		require.Equal(t, srv.monitor.CheckRT(1), false)
	}()
}
