package gproxy_client

import (
	"context"
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"strconv"
)

// Init gproxy client
func NewClient(cfg *config.Config, log *zap.SugaredLogger) gproxy.ProxyClient {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ProxyService.GRPCTimeout)
	defer cancel()
	target := ":" + strconv.Itoa(cfg.ProxyService.GRPCPort)
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
	logger.Panic(err, log, "grpc.DialContext error")
	return gproxy.NewProxyClient(conn)
}
