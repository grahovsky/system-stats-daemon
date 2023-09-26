package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os/signal"
	"syscall"

	pb "github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
	"github.com/grahovsky/system-stats-daemon/internal/client"
	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"google.golang.org/grpc"
)

var (
	host           string
	port           string
	responsePeriod int64
	rangeTime      int64
)

func init() {
	flag.StringVar(&host, "host", "0.0.0.0", "server host")
	flag.StringVar(&port, "port", "8086", "server port")
	flag.Int64Var(&responsePeriod, "n", 5, "period for sending statistics (sec)")
	flag.Int64Var(&rangeTime, "m", 15, "the range for which the average statistics are collected (sec)")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	conn, err := grpc.Dial(net.JoinHostPort(host, port), grpc.WithInsecure())
	if err != nil {
		fmt.Println("failed to dial connection" + err.Error())
		return
	}
	defer conn.Close()

	pbC := pb.NewStatsServiceClient(conn)

	req := &pb.StatsRequest{
		ResponsePeriod: responsePeriod,
		RangeTime:      rangeTime,
	}

	stream, err := pbC.StatsMonitoring(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	go func() {
		for {
			data, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Println("stats done..")
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
