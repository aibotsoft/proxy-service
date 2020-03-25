package migration

import (
	"github.com/stretchr/testify/assert"
	"proxy-service/internal/storage"
	"testing"
)

func TestUp(t *testing.T) {
	db, err := storage.NewStorage()
	assert.Nil(t, err)
	err = Up(db)
	assert.Nil(t, err)
}

func TestUpTo(t *testing.T) {
	targetVersion := 1
	db, err := storage.NewStorage()
	assert.Nil(t, err)
	err = UpTo(db, targetVersion)
	assert.Nil(t, err)
}
