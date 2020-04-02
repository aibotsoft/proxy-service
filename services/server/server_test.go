package server_test

import (
	"context"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/aibotsoft/proxy-service/pkg/utils"
	"github.com/aibotsoft/proxy-service/services/server"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	utils.MustLoadDevEnv("../../.env")
}

func initServer(t *testing.T) *server.Server {
	t.Helper()
	log := logger.New()
	cfg := config.New()
	db := sqlserver.MustConnect(cfg)
	s := server.NewServer(cfg, log, db)
	return s
}

func TestNewServer(t *testing.T) {
	var err error
	cfg := config.New()
	log := logger.New()
	log.Info(cfg)
	db := sqlserver.MustConnect(cfg)
	defer db.Close()
	srv := server.NewServer(cfg, log, db)
	go func() {
		err = srv.Serve()
	}()
	time.Sleep(time.Millisecond * 100)
	assert.NoError(t, err)
	srv.GracefulStop()
}
func TestServer_GetNextProxy(t *testing.T) {
	log := logger.New()
	cfg := config.New()
	db := sqlserver.MustConnect(cfg)
	s := server.NewServer(cfg, log, db)
	got, err := s.GetNextProxy(context.Background(), nil)
	if assert.NoError(t, err, err) {
		assert.IsType(t, int64(0), got.GetProxyItem().GetProxyId())
	}
}
func TestServer_GetBestProxy(t *testing.T) {
	s := initServer(t)
	got, err := s.GetBestProxy(context.Background(), nil)
	if assert.NoError(t, err, err) {
		assert.IsType(t, int64(0), got.GetProxyItem().GetProxyId())
	}
}
