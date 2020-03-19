package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	err := LoadEnv()
	assert.Nil(t, err)
	assert.Equal(t, "true", os.Getenv("TEST_LOAD_ENV"))
}

func TestNew(t *testing.T) {
	config, err := New()
	assert.Nil(t, err)
	assert.Equal(t, config.Database, "")
}
