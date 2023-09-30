package tests

import (
	"context"
	"net"
	"testing"
	"time"

	pb "github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
	"github.com/grahovsky/system-stats-daemon/internal/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestGRPCServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	srv := service.NewStatsMonitoringSever(ctx)

	grpcSrv := grpc.NewServer()
	pb.RegisterStatsServiceServer(grpcSrv, srv)

	dialer := func() func(context.Context, string) (net.Conn, error) {
		lis := bufconn.Listen(1024 * 1024)

		go func() {
			if err := grpcSrv.Serve(lis); err != nil {
				require.NoError(t, err)
			}
		}()

		return func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}
	}

	getConn := func(ctx context.Context) *grpc.ClientConn {
		conn, err := grpc.DialContext(ctx, "",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(dialer()))
		if err != nil {
			t.Fatal(err)
		}
		return conn
	}

	t.Run("get stats client", func(t *testing.T) {
		conn := getConn(ctx)
		defer conn.Close()

		pbC := pb.NewStatsServiceClient(conn)
		req := &pb.StatsRequest{
			ResponsePeriod: 1,
			RangeTime:      1,
		}

		stream, err := pbC.StatsMonitoring(ctx, req)
		require.NoError(t, err)

		resp, err := stream.Recv()
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}
