package client

import "github.com/spf13/pflag"

type Config struct {
	Host           string
	Port           string
	ResponsePeriod int64
	RangeTime      int64
}

func ParseConfig(c *Config) {
	pflag.StringVarP(&c.Host, "host", "h", "localhost", "server host")
	pflag.StringVarP(&c.Port, "port", "p", "8086", "server port")
	pflag.Int64VarP(&c.ResponsePeriod, "responseperiod", "n", 5, "period for sending statistics (sec)")
	pflag.Int64VarP(&c.RangeTime, "rangetime", "m", 15, "the range for which the average statistics are collected (sec)")
	pflag.Parse()
}
