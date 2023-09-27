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

	srv := NewStatsMonitoringSever(ctx)
	go func() {
		srv.StartMonitoring()
		require.Equal(t, srv.checkRT(1), false)
	}()
}
