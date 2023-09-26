package service

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
	"github.com/grahovsky/system-stats-daemon/internal/config"
	"google.golang.org/grpc"
)

type StatsMonitoringSever struct {
	ctx        context.Context
	cStorage   *cStorage
	grpcServer *grpc.Server
	pb.UnimplementedStatsServiceServer
}

func NewStatsMonitoringSever(ctx context.Context) *StatsMonitoringSever {
	return &StatsMonitoringSever{
		ctx: ctx,
	}
}

func (s *StatsMonitoringSever) Start() error {
	lis, err := net.Listen("tcp", net.JoinHostPort(config.Settings.Server.Host, config.Settings.Server.Port))
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	s.grpcServer = grpc.NewServer(opts...)
	defer s.grpcServer.GracefulStop()

	s.StartMonitoring()
	pb.RegisterStatsServiceServer(s.grpcServer, s)

	if err = s.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *StatsMonitoringSever) Stop() {
	s.grpcServer.GracefulStop()
}

func (s *StatsMonitoringSever) StatsMonitoring(req *pb.StatsRequest, stream pb.StatsService_StatsMonitoringServer) error {
	tiker := time.NewTicker(5 * time.Second)
	defer tiker.Stop()

	for {
		select {
		case <-tiker.C:
			fmt.Println("grpc tik")
		case <-s.ctx.Done():
			fmt.Println("grpc done")
			return nil
		}
	}
}
