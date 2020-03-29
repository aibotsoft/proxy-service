package check

import (
	"github.com/aibotsoft/gproxy"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_CheckAddr(t *testing.T) {
	cfg := config.New()
	log := logger.New()
	c := New(cfg, log, nil)
	t.Run("Fail in localhost", func(t *testing.T) {
		c.cfg.ProxyService.CheckTimeout = 10 * time.Millisecond
		p := &gproxy.ProxyItem{
			ProxyIp:   "0.0.0.0",
			ProxyPort: 80,
		}
		got := c.CheckProxy(p)
		assert.False(t, got.ConnStatus)
	})
	t.Run("Suc in 1.1.1.1", func(t *testing.T) {
		c.cfg.ProxyService.CheckTimeout = 100 * time.Millisecond
		p := &gproxy.ProxyItem{
			ProxyIp:   "1.1.1.1",
			ProxyPort: 80,
		}
		got := c.CheckProxy(p)
		assert.True(t, got.ConnStatus)
	})
}
