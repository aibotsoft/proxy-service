package collect

import (
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/proxy-service/internal/utils"
	"github.com/antchfx/htmlquery"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"strconv"
)

const (
	proxyTablePath   = `//table[@id="proxylisttable"]/tbody/tr`
	proxyIP          = 0
	proxyPort        = 1
	proxyCountryCode = 2
	proxyCountryName = 3
	proxyAnonymity   = 4
)

func (c *Collect) collectProxy() ([]gproxy.ProxyItem, error) {
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

func (c *Collect) scrapeProxy(proxyPageNode *html.Node) ([]gproxy.ProxyItem, error) {
	tr, err := htmlquery.QueryAll(proxyPageNode, proxyTablePath)
	if err != nil {
		return nil, errors.Wrap(err, "htmlquery.QueryAll error")
	}
	var proxyList []gproxy.ProxyItem
	for _, row := range tr {
		td, err := htmlquery.QueryAll(row, "/td")
		if err != nil {
			c.log.Error("Error QueryAll: ", err)
			continue
		}
		port, err := strconv.Atoi(htmlquery.InnerText(td[proxyPort]))
		if err != nil {
			c.log.Error("proxyPort to int convert error: ", err)
			continue
		}
		item := gproxy.ProxyItem{
			ProxyIp:   htmlquery.InnerText(td[proxyIP]),
			ProxyPort: int64(port),
			ProxyCountry: &gproxy.ProxyCountry{
				CountryName: htmlquery.InnerText(td[proxyCountryName]),
				CountryCode: htmlquery.InnerText(td[proxyCountryCode]),
			},
		}
		proxyList = append(proxyList, item)
	}
	return proxyList, nil
}
