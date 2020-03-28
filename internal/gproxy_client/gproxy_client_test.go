package gproxy_client

import (
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	cfg := config.New()
	log := logger.New()
	c := NewClient(cfg, log)
	assert.NotEmpty(t, c)

}
