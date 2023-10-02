package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	pb "github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
	"github.com/grahovsky/system-stats-daemon/internal/config"
	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/grahovsky/system-stats-daemon/internal/monitor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var (
	ErrResponsePeriod = errors.New("incorrect response period")
	ErrRangeTime      = errors.New("incorrect range time")
)

type StatsMonitoringSever struct {
	ctx        context.Context
	monitor    *monitor.Server
	grpcServer *grpc.Server
	pb.UnimplementedStatsServiceServer
}

func NewStatsMonitoringSever(ctx context.Context) *StatsMonitoringSever {
	monitor := monitor.NewMonitor(ctx)
	monitor.StartMonitoring()

	return &StatsMonitoringSever{
		ctx:     ctx,
		monitor: monitor,
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
		return status.Error(codes.Internal, fmt.Errorf("rangetime (%d) exceeds limit stored stats time (%d) %w",
			req.RangeTime, config.Settings.Stats.Limit, ErrRangeTime).Error())
	}

	if req.RangeTime < 1 {
		logger.Error(fmt.Sprintf("client %s err, min rangetime 1 sec", clientID))
		return status.Error(codes.Internal, fmt.Errorf("rangetime (%d). min rangetime 1 sec %w",
			req.RangeTime, ErrRangeTime).Error())
	}

	if req.ResponsePeriod < 1 {
		logger.Error(fmt.Sprintf("client %s err, min response period 1 sec", clientID))
		return status.Error(codes.Internal, fmt.Errorf("response period (%d). min response period 1 sec %w",
			req.ResponsePeriod, ErrResponsePeriod).Error())
	}

	responseTicker := time.NewTicker(time.Duration(req.ResponsePeriod) * time.Second)
	defer responseTicker.Stop()
	for {
		select {
		case <-responseTicker.C:
			if !s.monitor.CheckRT(req.RangeTime) {
				continue
			}
			err := stream.Send(&pb.StatsResponse{
				LoadInfo:    s.monitor.LoadInfoAvg(req.RangeTime),
				CpuInfo:     s.monitor.CPUInfoAvg(req.RangeTime),
				DiskInfo:    s.monitor.DiskInfoAvg(req.RangeTime),
				TalkersInfo: s.monitor.TalkersAvg(req.RangeTime),
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
