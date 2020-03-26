package main

import (
	"fmt"
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/postgres"
	"github.com/aibotsoft/proxy-service/cmd/collect"
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
	db, err := postgres.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
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
	collectService := collect.New(cfg, log)
	collectService.Start()
	log.Info("exit: ", <-errc)
	func() {
		s.GracefulStop()
		collectService.Stop()
	}()

	//msg, err := msg_server.NewMsgServer(cfg, log)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Fatal(msg.Run())
	//log.Print(ec)
	//log.Print(cfg.Controller.NewProxyAddress)
	//_, err = ec.Subscribe(cfg.Controller.NewProxyAddress, newProxyHandler)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//db, err := storage.Connect(cfg)
	//if err != nil {
	//	panic(err)
	//}

	//for {
	//	log.Print(time.Now())
	//	time.Sleep(time.Second * 10)
	//}
	//err = migration.Up(db)
	//if err != nil {
	//	panic(err)
	//}
	//log.Println("main : Started : Initializing API support")
	//server := http.Server{
	//	Addr:         cfg.Web.APIHost,
	//	Handler:      handlers.API(),
	//	ReadTimeout:  cfg.Web.ReadTimeout,
	//	WriteTimeout: cfg.Web.WriteTimeout,
	//}
	//log.Fatal(server.ListenAndServe())

}
