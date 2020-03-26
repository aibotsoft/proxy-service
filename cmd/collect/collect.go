package collect

import (
	"context"
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/micro/config"
	"github.com/dgraph-io/ristretto"
	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
	"time"
)

const (
	gPort = ":50051"
)

type Collect struct {
	cfg         *config.Config
	cron        *cron.Cron
	log         *zap.SugaredLogger
	client      *resty.Client
	proxyClient gproxy.ProxyClient
	cache       *ristretto.Cache
}

func (c *Collect) Start() {
	c.cron.Schedule(cron.Every(c.cfg.ProxyService.CollectPeriod), cron.FuncJob(c.CollectJob))
	c.cron.Start()
}
func (c *Collect) Stop() {
	c.cron.Stop()
}

func (c *Collect) CollectJob() {
	proxyItems, err := c.collectProxy()
	if err != nil {
		c.log.Error("collectProxy error: ", err)
		return
	}
	for _, p := range proxyItems {
		c.sendProxyItem(&p)
	}
	c.log.Debug(c.cache.Metrics)
}
func (c *Collect) proxyItemKey(p *gproxy.ProxyItem) string {
	return p.GetProxyIp() + ":" + strconv.Itoa(int(p.GetProxyPort()))
}
func (c *Collect) sendProxyItem(p *gproxy.ProxyItem) {
	_, ok := c.cache.Get(c.proxyItemKey(p))
	if ok {
		return
	}
	c.log.Debug(p)
	req := &gproxy.CreateProxyRequest{ProxyItem: p}
	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.ProxyService.GRPCTimeout)
	defer cancel()

	res, err := c.proxyClient.CreateProxy(ctx, req)
	if err != nil {
		c.log.Error("proxyClient.CreateProxy error: ", err)
	}
	c.log.Debug(res.GetProxyItem().GetProxyId())
	ok = c.cache.Set(c.proxyItemKey(p), nil, 1)
}

func New(cfg *config.Config, log *zap.SugaredLogger) *Collect {
	c := cron.New()

	tr := &http.Transport{TLSHandshakeTimeout: 0 * time.Second}
	client := resty.New().SetTransport(tr).EnableTrace()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, gPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	proxyClient := gproxy.NewProxyClient(conn)

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,
		MaxCost:     1e5,
		BufferItems: 64,
		Metrics:     true,
	})
	if err != nil {
		log.Fatalf("did not get cache: %v", err)
	}

	return &Collect{
		cfg:         cfg,
		cron:        c,
		log:         log,
		client:      client,
		proxyClient: proxyClient,
		cache:       cache,
	}
}
