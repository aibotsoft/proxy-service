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
)

func initStore(t *testing.T) *server.Store {
	t.Helper()
	gotenv.Must(gotenv.Load, "./../../.env")
	log := logger.New()
	cfg := config.New()
	db := sqlserver.MustConnect(cfg)
	assert.NotEmpty(t, db, db)
	store := server.NewStore(cfg, log, db)
	assert.NotEmpty(t, store, store)
	return store
}
func TestStore(t *testing.T) {
	ctx := context.Background()
	store := initStore(t)

	t.Run("GetOrCreateProxyCountry", func(t *testing.T) {
		p := &pb.ProxyCountry{
			CountryName: "Unknown",
			CountryCode: "22",
		}
		err := store.GetOrCreateProxyCountry(ctx, p)
		if assert.NoError(t, err, err) {
			assert.Equal(t, int64(1), p.CountryId, p)
		}
	})
	t.Run("GetOrCreateProxyItem", func(t *testing.T) {
		p := &pb.ProxyItem{
			ProxyAddr: "0.0.0.0:80",
			ProxyCountry: &pb.ProxyCountry{
				CountryName: "Unknown",
				CountryCode: "NA",
			},
		}
		err := store.GetOrCreateProxyItem(ctx, p)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(1), p.ProxyId)
		}
		// repeat query for get from cache
		err = store.GetOrCreateProxyItem(ctx, p)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(1), p.ProxyId)
		}
	})
}

func TestStore_CreateProxyStat(t *testing.T) {
	store := initStore(t)
	stat := &pb.ProxyStat{
		ProxyId:    1,
		ConnTime:   1,
		ConnStatus: true,
	}
	err := store.CreateProxyStat(context.Background(), stat)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, stat.CreatedAt, stat.CreatedAt)
	}
}

func TestStore_GetNextProxyItemBatch(t *testing.T) {
	store := initStore(t)
	ctx := context.Background()
	err := store.GetNextProxyItemBatch(ctx, 100)
	assert.NoError(t, err)
}

func TestStore_GetNextProxyItem(t *testing.T) {
	store := initStore(t)
	get, err := store.GetNextProxyItem(context.Background())
	if assert.NoError(t, err) {
		assert.NotEmpty(t, get, get)
	}
}
