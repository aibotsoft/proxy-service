package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"proxy-service/internal/models"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()

}

func TestOpen(t *testing.T) {
	var v int
	db, err := Open()
	assert.Nil(t, err)
	err = db.QueryRow(context.Background(), "select 1").Scan(&v)
	assert.Nil(t, err)
	assert.Equal(t, 1, v)
}

func TestSelectProxy(t *testing.T) {
	query := "select proxy_port, proxy_ip from proxy where proxy_id=1"
	p := models.Proxy{}
	db, err := Open()
	assert.Nil(t, err)
	err = db.QueryRow(context.Background(), query).Scan(&p.Port, &p.Ip)
	assert.Nil(t, err)
	log.Print(p)
}
func TestSelectCountryId(t *testing.T) {
	query := "select find_country_by_code('NA')"
	c := models.Country{}
	db, err := Open()
	assert.Nil(t, err)
	err = db.QueryRow(context.Background(), query).Scan(&c.CountryId)
	assert.Nil(t, err)
	log.Print(c)
}
