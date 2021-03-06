package check

import (
	"context"
	pb "github.com/aibotsoft/gen/proxypb"
	"github.com/aibotsoft/micro/config"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

type Check struct {
	cfg         *config.Config
	log         *zap.SugaredLogger
	cron        *cron.Cron
	proxyClient pb.ProxyClient
}

func New(cfg *config.Config, log *zap.SugaredLogger, proxyClient pb.ProxyClient) *Check {
	return &Check{cfg: cfg, cron: cron.New(), log: log, proxyClient: proxyClient}
}
func (c *Check) Start() {
	c.cron.Schedule(cron.Every(c.cfg.ProxyService.CheckPeriod), cron.FuncJob(c.CheckProxyJob))
	c.cron.Schedule(cron.Every(c.cfg.ProxyService.DeleteBadProxyPeriod), cron.FuncJob(c.DeleteBadProxyJob))
	c.cron.Schedule(cron.Every(c.cfg.ProxyService.DeleteOldStatPeriod), cron.FuncJob(c.DeleteOldStatJob))
	c.cron.Start()
}

func (c *Check) Stop() context.Context {
	ctx := c.cron.Stop()
	return ctx
}

// Регулярно помечаем удаленными (deleted_at) прокси с очень плохой статистикой
func (c *Check) DeleteOldStatJob() {
	c.log.Debug("DeleteOldStatJob begins..")
	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.ProxyService.GrpcTimeout)
	defer cancel()
	res, err := c.proxyClient.DeleteOldStat(ctx, nil)
	if err != nil {
		c.log.Info(errors.Wrap(err, "proxyClient.DeleteOldStat error"))
		return
	}
	for _, p := range res.GetDeletedStat() {
		c.log.Infof("Deleted old stat %+v", p)
	}
	c.log.Debug("DeleteOldStatJob end ", res.GetDeletedStat())

}

// Регулярно помечаем удаленными (deleted_at) прокси с очень плохой статистикой
func (c *Check) DeleteBadProxyJob() {
	c.log.Info("DeleteBadProxyJob begins..")
	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.ProxyService.GrpcTimeout)
	defer cancel()
	res, err := c.proxyClient.DeleteBadProxy(ctx, nil)
	if err != nil {
		c.log.Info(errors.Wrap(err, "proxyClient.DeleteBadProxy error"))
		return
	}
	for _, p := range res.GetBadProxy() {
		c.log.Infof("Deleted a bad proxy %+v", p)
	}
	c.log.Debug("DeleteBadProxyJob end ", res.GetBadProxy())
}

// Основная работа по проверке прокси
func (c *Check) CheckProxyJob() {
	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.ProxyService.CheckTimeout*2)
	defer cancel()
	// Получаем следующее прокси для проверки
	pi, err := c.getNextProxyItem(ctx)
	if err != nil {
		return
	}

	// Проверяем это прокси
	proxyStat := c.checkProxy(ctx, pi)

	// Отправляем результат проверки на сервер
	_, err = c.sendProxyStat(ctx, proxyStat)
	if err != nil {
		c.log.Error(err)
	}
}

func (c *Check) getNextProxyItem(ctx context.Context) (*pb.ProxyItem, error) {
	res, err := c.proxyClient.GetNextProxy(ctx, &pb.GetNextProxyRequest{})
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
func (c *Check) checkProxy(ctx context.Context, p *pb.ProxyItem) *pb.ProxyStat {
	ConnTime, ConnStatus := c.checkAddr(ctx, p.ProxyAddr)
	stat := &pb.ProxyStat{
		ProxyId:    p.ProxyId,
		ConnTime:   ConnTime.Milliseconds(),
		ConnStatus: ConnStatus,
	}
	return stat
}

func (c *Check) sendProxyStat(ctx context.Context, stat *pb.ProxyStat) (*pb.ProxyStat, error) {
	res, err := c.proxyClient.CreateProxyStat(ctx, &pb.CreateProxyStatRequest{ProxyStat: stat})
	if err != nil {
		return nil, errors.Wrap(err, "proxyClient.CreateProxyStat error")
	}
	return res.GetProxyStat(), nil
}

func (c *Check) checkAddr(ctx context.Context, addr string) (time.Duration, bool) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp4", addr, c.cfg.ProxyService.CheckTimeout)
	if err != nil {
		return time.Since(start), false
	}
	defer conn.Close()
	return time.Since(start), true
}
