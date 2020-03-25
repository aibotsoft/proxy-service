package msg_server

import (
	"log"
	"proxy-service/internal/broker"
	"proxy-service/internal/config"
	"proxy-service/internal/models"
	"proxy-service/internal/storage"
)

//Создаем брокер сервер, который принимает новые входящие прокси и заносит в базу
// Поступающие прокси кешируем чтобы часто не дергать базу данных
type MsgServer struct {
	cfg   *config.Config
	log   *log.Logger
	ec    *broker.Broker
	store *storage.Storage
}

func (m MsgServer) newProxyHandler(p *models.ProxyItem) {
	//spew.Dump(p)

	err := m.store.GetOrCreateProxyItem(p)
	if err != nil {
		m.log.Println(err)
	}
}

func (m MsgServer) Run() error {
	_, err := m.ec.Subscribe(m.cfg.Controller.NewProxyAddress, m.newProxyHandler)
	if err != nil {
		return err
	}
	select {}
	//panic("implement me")
}

func NewMsgServer(cfg *config.Config, log *log.Logger) (*MsgServer, error) {
	ec, err := broker.NewBroker(cfg)
	if err != nil {
		return nil, err
	}
	store, err := storage.NewStorage(cfg, log)
	if err != nil {
		return nil, err
	}

	return &MsgServer{cfg: cfg, log: log, ec: ec, store: store}, nil
}
