package collect

import (
	"context"
	pb "github.com/aibotsoft/gen/proxypb"
	"github.com/aibotsoft/micro/cache"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/proxy-service/pkg/utils"
	"github.com/antchfx/htmlquery"
	"github.com/dgraph-io/ristretto"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"net"
	"net/http"
	"time"
)

const (
	proxyTablePath   = `//table[@id="proxylisttable"]/tbody/tr`
	proxyIP          = 0
	proxyPort        = 1
	proxyCountryCode = 2
	proxyCountryName = 3
	proxyAnonymity   = 4
)

type Collect struct {
	cfg         *config.Config
	cron        *cron.Cron
	log         *zap.SugaredLogger
	client      *resty.Client
	proxyClient pb.ProxyClient
	cache       *ristretto.Cache
}

func New(cfg *config.Config, log *zap.SugaredLogger, proxyClient pb.ProxyClient) *Collect {
	c := cron.New()

	tr := &http.Transport{TLSHandshakeTimeout: 0 * time.Second}
	client := resty.New().SetTransport(tr).EnableTrace()

	return &Collect{
		cfg:         cfg,
		cron:        c,
		log:         log,
		client:      client,
		proxyClient: proxyClient,
		cache:       cache.NewCache(cfg),
	}
}
func (c *Collect) Start() {
	c.cron.Schedule(cron.Every(c.cfg.ProxyService.CollectPeriod), cron.FuncJob(c.CollectJob))
	c.cron.Start()
}
func (c *Collect) Stop() {
	c.cron.Stop()
}
func (c *Collect) CollectJob() {
	newProxyCount := 0
	proxyItems, err := c.collectProxy()
	if err != nil {
		c.log.Error("collectProxy error: ", err)
		return
	}
	for _, p := range proxyItems {
		res := c.sendProxyItem(&p)
		newProxyCount = newProxyCount + res
	}
	c.log.Debug(c.cache.Metrics)
	c.log.Debugf("Got %v proxy, bun only %v new", len(proxyItems), newProxyCount)
}

func (c *Collect) sendProxyItem(p *pb.ProxyItem) int {
	_, ok := c.cache.Get(p.ProxyAddr)
	if ok {
		return 0
	}
	req := &pb.CreateProxyRequest{ProxyItem: p}
	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.ProxyService.GrpcTimeout)
	defer cancel()
	res, err := c.proxyClient.CreateProxy(ctx, req)
	if err != nil {
		c.log.Error("proxyClient.CreateProxy error: ", err)
		return 0
	}
	ok = c.cache.Set(p.ProxyAddr, res.GetProxyItem().GetProxyId(), 1)
	return 1
}

func (c *Collect) collectProxy() ([]pb.ProxyItem, error) {
	pageNode, err := c.getNewProxy()
	if err != nil {
		return nil, err
	}
	return c.scrapeProxy(pageNode)
}

// Получаем прокси на сайте и парсим документ в html.Node
func (c *Collect) getNewProxy() (*html.Node, error) {
	resp, err := c.client.R().SetDoNotParseResponse(true).Get(c.cfg.ProxyService.CollectUrl)
	if err != nil {
		return nil, errors.Wrap(err, "getNewProxy error")
	}
	if config.IsDev() {
		utils.LogRestyTrace(c.log, resp)
	}
	body := resp.RawBody()
	defer body.Close()
	return html.Parse(body)
}

func (c *Collect) scrapeProxy(proxyPageNode *html.Node) ([]pb.ProxyItem, error) {
	tr, err := htmlquery.QueryAll(proxyPageNode, proxyTablePath)
	if err != nil {
		return nil, errors.Wrap(err, "htmlquery.QueryAll error")
	}
	var proxyList []pb.ProxyItem
	for _, row := range tr {
		td, err := htmlquery.QueryAll(row, "/td")
		if err != nil {
			c.log.Error("Error QueryAll: ", err)
			continue
		}
		item := pb.ProxyItem{
			ProxyAddr: net.JoinHostPort(htmlquery.InnerText(td[proxyIP]), htmlquery.InnerText(td[proxyPort])),
			ProxyCountry: &pb.ProxyCountry{
				CountryName: htmlquery.InnerText(td[proxyCountryName]),
				CountryCode: htmlquery.InnerText(td[proxyCountryCode]),
			},
		}
		proxyList = append(proxyList, item)
	}
	return proxyList, nil
}
