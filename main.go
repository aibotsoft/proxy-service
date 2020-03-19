package main

import (
	"fmt"
	"proxy-service/internal/broker"
	"proxy-service/internal/config"
	"proxy-service/internal/logging"
	"proxy-service/internal/models"
	"time"
)

func newProxyHandler(p *models.ProxyItem) {
	fmt.Printf("%+v", p)
}

func main() {
	cfg := config.NewConfig()
	log := logging.New(cfg)
	log.Println("Beginning...")
	log.Printf("Config: %+v", cfg)
	ec, err := broker.NewBroker(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(ec)
	log.Print(cfg.Controller.NewProxyAddress)
	_, err = ec.Subscribe(cfg.Controller.NewProxyAddress, newProxyHandler)
	if err != nil {
		log.Fatal(err)
	}

	for {
		log.Print(time.Now())
		time.Sleep(time.Second * 10)
	}
	//db, err := storage.Open()
	//if err != nil {
	//	panic(err)
	//}
	//err = db.Ping(context.Background())
	//if err != nil {
	//	panic(err)
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
