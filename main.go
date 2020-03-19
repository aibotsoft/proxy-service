package main

import (
	"proxy-service/internal/config"
	"proxy-service/internal/logging"
)

func main() {
	cfg := config.NewConfig()
	log := logging.New(cfg)
	log.Println("Beginning...")
	log.Printf("Config: %+v", cfg)
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
