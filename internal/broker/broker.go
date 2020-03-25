package broker

import (
	"github.com/nats-io/nats.go"
	"proxy-service/internal/config"
)

type Broker struct {
	cfg *config.Config
	*nats.EncodedConn
}

func NewBroker(cfg *config.Config) (*Broker, error) {
	//b := &Broker{cfg: cfg}
	opts := nats.Options{
		Url:            cfg.Broker.Url,
		AllowReconnect: cfg.Broker.AllowReconnect,
		MaxReconnect:   cfg.Broker.MaxReconnect,
		ReconnectWait:  cfg.Broker.ReconnectWait,
		Timeout:        cfg.Broker.Timeout,
	}
	nc, err := opts.Connect()
	if err != nil {
		return nil, err
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	//b.EncodedConn = ec
	return &Broker{cfg: cfg, EncodedConn: ec}, nil

}
