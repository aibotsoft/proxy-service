package main

import (
	"proxy-service/internal/config"
	"proxy-service/internal/logging"
	"proxy-service/internal/msg_server"
)

func main() {
	cfg := config.NewConfig()
	log := logging.New(cfg)
	log.Println("Beginning...")
	log.Printf("Config: %+v", cfg)

	msg, err := msg_server.NewMsgServer(cfg, log)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(msg.Run())

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
