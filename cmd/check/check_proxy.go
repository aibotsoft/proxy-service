package check

import (
	"github.com/aibotsoft/gproxy"
	"net"
	"strconv"
	"time"
)

func (c *Check) checkAddr(addr string) (time.Duration, bool) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp4", addr, c.cfg.ProxyService.CheckTimeout)
	if err != nil {
		//c.log.Error(err)
		return time.Since(start), false
	}
	defer conn.Close()
	return time.Since(start), true
}

func (c *Check) CheckProxy(p *gproxy.ProxyItem) *gproxy.ProxyStat {
	addr := net.JoinHostPort(p.ProxyIp, strconv.Itoa(int(p.ProxyPort)))
	ConnTime, ConnStatus := c.checkAddr(addr)
	stat := &gproxy.ProxyStat{
		ProxyId:    p.ProxyId,
		ConnTime:   ConnTime.Milliseconds(),
		ConnStatus: ConnStatus,
	}
	return stat
}
