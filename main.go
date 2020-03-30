package main

import (
	"fmt"
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/postgres"
	"github.com/aibotsoft/proxy-service/cmd/check"
	"github.com/aibotsoft/proxy-service/cmd/collect"
	"github.com/aibotsoft/proxy-service/internal/gproxy_client"
	"github.com/subosito/gotenv"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Init dependencies
	gotenv.Must(gotenv.Load)
	cfg := config.New()
	log := logger.New().With("service", cfg.Service.Name)
	log.Infow("Begin service", "config", cfg)
	db := postgres.MustConnect(cfg)
	defer db.Close()

	// Инициализируем GracefulStop
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()
	// Run gPRC proxy server
	s := gproxy.NewServer(db)
	go func() {
		errc <- s.Serve()
	}()
	// Init gProxy Client
	proxyClient := gproxy_client.NewClient(cfg, log)

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
	//msg, err := msg_server.NewMsgServer(cfg, log)
	//err = migration.Up(db)
	//log.Println("main : Started : Initializing API support")
	//server := http.Server{
	//	Addr:         cfg.Web.APIHost,
	//	Handler:      handlers.API(),
	//	ReadTimeout:  cfg.Web.ReadTimeout,
	//	WriteTimeout: cfg.Web.WriteTimeout,
	//}
	//log.Fatal(server.ListenAndServe())

}
