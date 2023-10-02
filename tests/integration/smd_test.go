// go:build integration

package integration_test

import (
	"context"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
	"github.com/grahovsky/system-stats-daemon/internal/client"
	"github.com/grahovsky/system-stats-daemon/internal/service"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SystemStatsSuite struct {
	suite.Suite
	statConn   *grpc.ClientConn
	statClient stats_service.StatsServiceClient
	ctx        context.Context
	respEmpty  *stats_service.StatsResponse
}

func (s *SystemStatsSuite) SetupSuite() {
	var cfg client.Config
	client.ParseConfig(&cfg)

	statConn, err := grpc.Dial(net.JoinHostPort(cfg.Host, cfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.statConn = statConn
	s.statClient = stats_service.NewStatsServiceClient(s.statConn)

	// s.ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	ctx := context.Background()
	s.ctx = ctx

	respEmpty := stats_service.StatsResponse{}
	s.respEmpty = &respEmpty
}

func (s *SystemStatsSuite) SetupTest() {
	seed := time.Now().UnixNano()
	rand.NewSource(seed)
	s.T().Log("seed:", seed)
}

func (s *SystemStatsSuite) TearDownSuite() {
	err := s.statConn.Close()
	if err != nil {
		s.T().Log(err)
	}
}

func (s *SystemStatsSuite) TearDownTest() {
	s.T().Logf("%s - done", s.T().Name())
}

func (s *SystemStatsSuite) TestStandard() {
	req := stats_service.StatsRequest{ResponsePeriod: 1, RangeTime: 1}

	stream, err := s.statClient.StatsMonitoring(s.ctx, &req)
	s.Require().NoError(err)

	resp, err := stream.Recv()
	s.Require().NoError(err)
	s.Require().NotEqual(resp, s.respEmpty)
}

func (s *SystemStatsSuite) TestResponsePeriod() {
	req := stats_service.StatsRequest{ResponsePeriod: 2, RangeTime: 2}

	stream, err := s.statClient.StatsMonitoring(s.ctx, &req)
	s.Require().NoError(err)

	time1 := time.Now()
	_, err = stream.Recv()
	time2 := time.Now()
	s.Require().NoError(err)
	s.Require().Equal(2, int(time2.Sub(time1).Seconds()))
}

func (s *SystemStatsSuite) TestErrorResponsePeriod() {
	req := stats_service.StatsRequest{ResponsePeriod: -10, RangeTime: 1}

	stream, err := s.statClient.StatsMonitoring(s.ctx, &req)
	s.Require().NoError(err)

	_, err = stream.Recv()
	s.Require().Error(err)
	s.Require().Error(service.ErrResponsePeriod, err)
}

func (s *SystemStatsSuite) TestErrorRangeTime() {
	req := stats_service.StatsRequest{ResponsePeriod: 1, RangeTime: 100000}

	stream, err := s.statClient.StatsMonitoring(s.ctx, &req)
	s.Require().NoError(err)

	_, err = stream.Recv()
	s.Require().Error(err)
	s.Require().Error(service.ErrRangeTime, err)
}

func TestSystemStatusSuite(t *testing.T) {
	suite.Run(t, new(SystemStatsSuite))
}
