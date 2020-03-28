package check

import (
	"context"
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/micro/config"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Check struct {
	cfg         *config.Config
	log         *zap.SugaredLogger
	cron        *cron.Cron
	proxyClient gproxy.ProxyClient
}

func (c *Check) getNextProxyItem() (*gproxy.ProxyItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.ProxyService.GRPCTimeout)
	defer cancel()
	res, err := c.proxyClient.GetNextProxy(ctx, &gproxy.GetNextProxyRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "proxyClient.GetNextProxy error")
	}
	return res.GetProxyItem(), nil
}

func New(cfg *config.Config, log *zap.SugaredLogger, cron *cron.Cron, proxyClient gproxy.ProxyClient) *Check {
	return &Check{cfg: cfg, cron: cron, log: log, proxyClient: proxyClient}
}

func (c *Check) Start() {
	//c.cron.Schedule(cron.Every(c.cfg.ProxyService.CollectPeriod), cron.FuncJob(c.CollectJob))
	c.cron.Start()
}
func (c *Check) Stop() {
	c.cron.Stop()
}
