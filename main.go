package main

import (
	"context"
	"fmt"
	"github.com/aibotsoft/gen/proxypb"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"time"

	"github.com/aibotsoft/proxy-service/pkg/utils"
	"github.com/aibotsoft/proxy-service/services/check"
	"github.com/aibotsoft/proxy-service/services/collect"
	"github.com/aibotsoft/proxy-service/services/server"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	// Init dependencies
	utils.MustLoadDevEnv()
	cfg := config.New()
	log := logger.New().With("service", cfg.Service.Name)
	log.Infow("Begin service", "config", cfg)
	//db := postgres.MustConnect(cfg)
	db := sqlserver.MustConnect(cfg)
	defer db.Close()

	// Инициализируем GracefulStop
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()
	// Run gPRC proxy server
	s := server.NewServer(cfg, log, db)
	go func() {
		errc <- s.Serve()
	}()
	// Init gProxy Client
	//proxyClient := gproxy_client.NewClient(cfg, log)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, net.JoinHostPort("", strconv.Itoa(cfg.ProxyService.GrpcPort)), grpc.WithInsecure(), grpc.WithBlock())
	logger.Panic(err, log, "grpc.DialContext error")
	proxyClient := proxypb.NewProxyClient(conn)

	// Run Collect Service
	collectService := collect.New(cfg, log, proxyClient)
	collectService.Start()

	// Run Check Service
	checkService := check.New(cfg, log, proxyClient)
	checkService.Start()

	// Closing services
	defer func() {
		log.Debug("begin closing services")
		s.GracefulStop()
		collectService.Stop()
		checkService.Stop()
	}()

	log.Info("exit: ", <-errc)
	//err = migration.Up(db)
}
