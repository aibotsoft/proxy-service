package check

import (
	"context"
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/micro/config"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Check struct {
	cfg         *config.Config
	log         *zap.SugaredLogger
	cron        *cron.Cron
	proxyClient gproxy.ProxyClient
}

func (c *Check) CheckProxyJob() {
	pi, err := c.getNextProxyItem()
	if err != nil {
		return
	}
	proxyStat := c.CheckProxy(pi)
	_, err = c.SendProxyStat(proxyStat)
	if err != nil {
		c.log.Error(err)
	}
}
func (c *Check) getNextProxyItem() (*gproxy.ProxyItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.ProxyService.GRPCTimeout)
	defer cancel()
	res, err := c.proxyClient.GetNextProxy(ctx, &gproxy.GetNextProxyRequest{})
	switch {
	case status.Convert(err).Code() == codes.NotFound:
		c.log.Info(err)
		return nil, errors.Wrap(err, "proxyClient.GetNextProxy error")
	case err != nil:
		c.log.Error(err)
		return nil, errors.Wrap(err, "proxyClient.GetNextProxy error")
	}
	return res.GetProxyItem(), nil
}

func (c *Check) SendProxyStat(stat *gproxy.ProxyStat) (*gproxy.ProxyStat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.ProxyService.GRPCTimeout)
	defer cancel()
	res, err := c.proxyClient.CreateProxyStat(ctx, &gproxy.CreateProxyStatRequest{ProxyStat: stat})
	if err != nil {
		return nil, errors.Wrap(err, "proxyClient.CreateProxyStat error")
	}
	return res.GetProxyStat(), nil
}

func New(cfg *config.Config, log *zap.SugaredLogger, proxyClient gproxy.ProxyClient) *Check {
	return &Check{cfg: cfg, cron: cron.New(), log: log, proxyClient: proxyClient}
}
func (c *Check) Start() {
	c.cron.Schedule(cron.Every(c.cfg.ProxyService.CheckPeriod), cron.FuncJob(c.CheckProxyJob))
	c.cron.Start()
}

func (c *Check) Stop() context.Context {
	ctx := c.cron.Stop()
	return ctx
}
