package service

import (
	"context"
	"net"
	"time"

	pb "github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
	"github.com/grahovsky/system-stats-daemon/internal/config"
	"github.com/grahovsky/system-stats-daemon/internal/logger"
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
	responseTicker := time.NewTicker(time.Duration(req.ResponsePeriod) * time.Second)
	defer responseTicker.Stop()
	for {
		select {
		case <-responseTicker.C:
			err := stream.Send(&pb.StatsResponse{
				LoadInfo: s.LoadInfoAvg(req.RangeTime),
				CpuInfo:  s.CpuInfoAvg(req.RangeTime),
				DiskInfo: s.DiskInfoAvg(req.RangeTime),
			})
			if err != nil {
				return err
			}
		case <-stream.Context().Done():
			logger.Error("sending data interrupted")
			return nil
		}
	}
}
