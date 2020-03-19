package storage

import (
	"context"
	"github.com/jackc/pgx/v4"
	"proxy-service/internal/config"
)

func Open() (*pgx.Conn, error) {
	config.MustLoadEnv()
	connConfig, err := pgx.ParseConfig("")
	if err != nil {
		return nil, err
	}
	//log.Print("connConfig: ", connConfig)
	return pgx.ConnectConfig(context.Background(), connConfig)
	//return pgx.Connect(context.Background(), "")
}
