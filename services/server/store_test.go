package server_test

import (
	"context"
	pb "github.com/aibotsoft/gen/proxypb"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/aibotsoft/proxy-service/services/server"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	gotenv.Must(gotenv.Load, "./../../.env")
	log := logger.New()
	cfg := config.New()
	db := sqlserver.MustConnect(cfg)
	assert.NotEmpty(t, db, db)
	store := server.NewStore(cfg, log, db)
	assert.NotEmpty(t, store, store)
	ctx := context.Background()
	//log.Info(store)

	t.Run("GetOrCreateProxyCountry", func(t *testing.T) {
		p := &pb.ProxyCountry{
			CountryName: "Unknown",
			CountryCode: "NA",
		}
		err := store.GetOrCreateProxyCountry(ctx, p)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(1), p.CountryId)
		}
		time.Sleep(time.Millisecond * 10)
		// repeat query for get from cache
		err = store.GetOrCreateProxyCountry(ctx, p)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(1), p.CountryId)
		}
	})
	t.Run("GetOrCreateProxyItem", func(t *testing.T) {
		p := &pb.ProxyItem{
			ProxyIp:   "0.0.0.0",
			ProxyPort: 80,
			ProxyCountry: &pb.ProxyCountry{
				CountryName: "Unknown",
				CountryCode: "NA",
			},
		}
		err := store.GetOrCreateProxyItem(ctx, p)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(1), p.ProxyId)
		}
	})

}
