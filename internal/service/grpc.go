package service

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
	"github.com/grahovsky/system-stats-daemon/internal/config"
	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
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
	addr := net.JoinHostPort(config.Settings.Server.Host, config.Settings.Server.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("started server %s", addr))

	var opts []grpc.ServerOption
	s.grpcServer = grpc.NewServer(opts...)
	defer s.grpcServer.GracefulStop()

	s.StartMonitoring()
	pb.RegisterStatsServiceServer(s.grpcServer, s)

	return s.grpcServer.Serve(lis)
}

func (s *StatsMonitoringSever) Stop() {
	s.grpcServer.GracefulStop()
}

func (s *StatsMonitoringSever) StatsMonitoring(req *pb.StatsRequest,
	stream pb.StatsService_StatsMonitoringServer,
) error {
	var clientID string
	if info, ok := peer.FromContext(stream.Context()); ok {
		clientID = info.Addr.String()
		logger.Info(fmt.Sprintf("connected client %s, rangetime: %d, responseperiod: %d",
			clientID, req.RangeTime, req.ResponsePeriod))
	}

	if req.RangeTime > config.Settings.Stats.Limit {
		logger.Error(fmt.Sprintf("client %s err, rangetime exceeds limit stored stats time", clientID))
		return status.Error(codes.Internal, fmt.Sprintf("rangetime (%d) exceeds limit stored stats time (%d)",
			req.RangeTime, config.Settings.Stats.Limit))
	}

	responseTicker := time.NewTicker(time.Duration(req.ResponsePeriod) * time.Second)
	defer responseTicker.Stop()
	for {
		select {
		case <-responseTicker.C:
			if !s.checkRT(req.RangeTime) {
				continue
			}
			err := stream.Send(&pb.StatsResponse{
				LoadInfo: s.LoadInfoAvg(req.RangeTime),
				CpuInfo:  s.CPUInfoAvg(req.RangeTime),
				DiskInfo: s.DiskInfoAvg(req.RangeTime),
			})
			if err != nil {
				return err
			}
		case <-s.ctx.Done():
			logger.Info("stopped server..")
			return nil
		case <-stream.Context().Done():
			logger.Error(fmt.Sprintf("client (%s), sending data interrupted", clientID))
			return nil
		}
	}
}
