package check

import (
	"context"
	"github.com/aibotsoft/gen/proxypb"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/aibotsoft/proxy-service/pkg/utils"
	"github.com/aibotsoft/proxy-service/services/server"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"testing"
	"time"
)

//func Test_CheckAddr(t *testing.T) {
//	cfg := config.New()
//	log := logger.New()
//	c := New(cfg, log, nil)
//	t.Run("Fail in localhost", func(t *testing.T) {
//		c.cfg.ProxyService.CheckTimeout = 10 * time.Millisecond
//		p := &gproxy.ProxyItem{
//			ProxyIp:   "0.0.0.0",
//			ProxyPort: 80,
//		}
//		got := c.checkProxy(p)
//		assert.False(t, got.ConnStatus)
//	})
//	t.Run("Suc in 1.1.1.1", func(t *testing.T) {
//		c.cfg.ProxyService.CheckTimeout = 100 * time.Millisecond
//		p := &gproxy.ProxyItem{
//			ProxyIp:   "1.1.1.1",
//			ProxyPort: 80,
//		}
//		got := c.checkProxy(p)
//		assert.True(t, got.ConnStatus)
//	})
//}
func initCheckService(t *testing.T) (*Check, func()) {
	t.Helper()
	utils.MustLoadDevEnv("./../../.env")

	var err error
	cfg := config.New()
	log := logger.New()
	log.Info(cfg)
	db := sqlserver.MustConnect(cfg)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	srv := server.NewServer(cfg, log, db)
	go func() {
		err = srv.Serve()
	}()
	time.Sleep(time.Millisecond * 100)
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, net.JoinHostPort("", strconv.Itoa(cfg.ProxyService.GrpcPort)), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	proxyClient := proxypb.NewProxyClient(conn)
	//srv.GracefulStop()
	checkService := New(cfg, log, proxyClient)

	closer := func() {
		log.Info("begin close all services")
		db.Close()
		srv.GracefulStop()
		err2 := conn.Close()
		if err != nil {
			log.Info(err2)
		}
		checkService.Stop()
	}
	return checkService, closer
}
func TestCheck_DeleteBadProxyJob(t *testing.T) {
	c, closer := initCheckService(t)
	defer closer()
	c.DeleteBadProxyJob()
}
func TestCheck_DeleteOldStatJob(t *testing.T) {
	c, closer := initCheckService(t)
	defer closer()
	c.DeleteOldStatJob()
}
