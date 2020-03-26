package collect

import (
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	cfg := config.New()
	log := logger.New()
	c := New(cfg, log)
	assert.NotEmpty(t, c.cfg)
	assert.NotEmpty(t, c.cron)
	assert.NotEmpty(t, c.client)
	assert.NotEmpty(t, c.proxyClient)
	assert.NotEmpty(t, c.cache)
}

func TestCollect_CollectJob(t *testing.T) {
	cfg := config.New()
	log := logger.New()
	c := New(cfg, log)
	c.CollectJob()
}
