package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os/signal"
	"syscall"

	pb "github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
	"github.com/grahovsky/system-stats-daemon/internal/client"
	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var cc client.ClientConfig

func init() {
	client.ParseConfig(&cc)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(net.JoinHostPort(cc.Host, cc.Port), opts...)
	if err != nil {
		fmt.Println("failed to dial connection" + err.Error())
		return
	}
	defer conn.Close()

	pbC := pb.NewStatsServiceClient(conn)

	req := &pb.StatsRequest{
		ResponsePeriod: cc.ResponsePeriod,
		RangeTime:      cc.RangeTime,
	}

	stream, err := pbC.StatsMonitoring(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	go func() {
		defer cancel()
		for {
			data, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Println("stats streaming done..")
				return
			}
			if err != nil {
				logger.Error(err.Error())
				return
			}
			client.PrintResponse(data)
		}
	}()

	<-ctx.Done()
}
