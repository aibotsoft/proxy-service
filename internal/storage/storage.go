package storage

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/jackc/pgx/v4"
	"log"
	"proxy-service/internal/config"
	"proxy-service/internal/models"
)

type Storage struct {
	cfg   *config.Config
	log   *log.Logger
	db    *pgx.Conn
	cache *ristretto.Cache
}

var selectProxyCountryIdByCode = `select country_id from country where country_code=$1`
var selectProxyIdByUI = `select proxy_id from proxy where proxy_ip=$1 and proxy_port=$2`

var insertProxyCountry = `INSERT INTO country (created_at, country_name, country_code) VALUES ($1, $2, $3) returning country_id`
var insertProxyItem = `INSERT INTO proxy_service.public.proxy (created_at, updated_at, proxy_ip, proxy_port, country_id) 
					VALUES ($1, $2, $3, $4, $5) returning proxy_id`

func (s Storage) GetOrCreateProxyItem(p *models.ProxyItem) error {
	err := s.SelectProxyIdByUI(p)
	if err == nil {
		return err
	}
	err = s.GetOrCreateProxyCountry(&p.Country)
	if err != nil {
		return err
	}
	return s.InsertProxyItem(p)
}

func (s Storage) InsertProxyItem(p *models.ProxyItem) error {
	return s.db.QueryRow(context.Background(), insertProxyItem, p.CreatedAt, p.UpdatedAt, p.ProxyIp, p.ProxyPort, p.Country.CountryId).Scan(&p.ProxyId)
}

func (s Storage) GetOrCreateProxyCountry(c *models.Country) error {
	err := s.SelectProxyCountryIdByCode(c)
	if err == pgx.ErrNoRows {
		return s.InsertProxyCountry(c)
	}
	return err
}

func (s Storage) SelectProxyIdByUI(p *models.ProxyItem) error {
	return s.db.QueryRow(context.Background(), selectProxyIdByUI, p.ProxyIp, p.ProxyPort).Scan(&p.ProxyId)
}

func (s Storage) SelectProxyCountryIdByCode(c *models.Country) error {
	return s.db.QueryRow(context.Background(), selectProxyCountryIdByCode, c.CountryCode).Scan(&c.CountryId)
}
func (s Storage) InsertProxyCountry(c *models.Country) error {
	return s.db.QueryRow(context.Background(), insertProxyCountry, c.CreatedAt, c.CountryName, c.CountryCode).Scan(&c.CountryId)
}

func NewStorage(cfg *config.Config, log *log.Logger) (*Storage, error) {
	ctx := context.Background()

	connConfig, err := pgx.ParseConfig("")
	if err != nil {
		return nil, err
	}
	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &Storage{cfg: cfg, log: log, db: conn}, nil
}
