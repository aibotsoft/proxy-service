package server

import (
	"context"
	"database/sql"
	pb "github.com/aibotsoft/gen/proxypb"
	"github.com/aibotsoft/micro/cache"
	"github.com/aibotsoft/micro/config"
	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

const (
	nextProxyQueryTimeout = 30 * time.Second
	//selectProxyCountryIdByCode = `select country_id from country where country_code=$1`
	//selectProxyIdByUI          = `select proxy_id from proxy where proxy_ip=$1 and proxy_port=$2`
	//getNextProxyItemBatch = `select proxy_id, proxy_ip, proxy_port from get_next_proxy_for_check(60, $1);`
	//insertProxyCountry = `INSERT INTO country (country_name, country_code) VALUES ($1, $2) returning country_id`
	//insertProxyStat = `INSERT INTO stat (proxy_id, conn_time, conn_status) VALUES ($1, $2, $3) returning created_at`
	//insertProxyItem = `INSERT INTO proxy_service.public.proxy (proxy_ip, proxy_port, country_id) VALUES ($1, $2, $3) returning proxy_id`
)

var ErrNoRows = errors.New("no rows in result set")

//var ErrQueryTooOften = errors.New("query too often, need to wait")

type Store struct {
	cfg                *config.Config
	log                *zap.SugaredLogger
	db                 *sql.DB
	cache              *ristretto.Cache
	nextProxyQueue     chan pb.ProxyItem
	nextProxyLastQuery time.Time
}

func NewStore(cfg *config.Config, log *zap.SugaredLogger, db *sql.DB) *Store {
	return &Store{log: log, db: db, cache: cache.NewCache(cfg),
		nextProxyQueue: make(chan pb.ProxyItem, 200)}
}

func (s *Store) GetOrCreateProxyCountry(ctx context.Context, c *pb.ProxyCountry) error {
	if get, b := s.cache.Get(c.CountryCode); b {
		if value, ok := get.(int64); ok {
			c.CountryId = value
			return nil
		}
	}
	err := s.db.QueryRowContext(ctx, "uspGetOrCreateProxyCountry",
		sql.Named("country_code", &c.CountryCode),
		sql.Named("country_name", &c.CountryName),
	).Scan(&c.CountryId)
	if err != nil {
		return errors.Wrap(err, "uspGetOrCreateProxyCountry error")
	}
	s.cache.Set(c.CountryCode, c.CountryId, 1)
	return nil
}

func (s *Store) GetOrCreateProxyItem(ctx context.Context, p *pb.ProxyItem) error {
	err := s.GetOrCreateProxyCountry(ctx, p.ProxyCountry)
	if err != nil {
		return errors.Wrap(err, "GetOrCreateProxyCountry error")
	}
	err = s.db.QueryRowContext(ctx, "uspGetOrCreateProxy",
		sql.Named("proxy_ip", &p.ProxyIp),
		sql.Named("proxy_port", &p.ProxyPort),
		sql.Named("country_id", &p.ProxyCountry.CountryId),
	).Scan(&p.ProxyId)
	if err != nil {
		return errors.Wrap(err, "uspGetOrCreateProxy error")
	}
	s.cache.Set(p.P, p.ProxyId, 1)

	return nil
}

//func (s *Store) SelectProxyCountryIdByCode(c *pb.ProxyCountry) error {
//	return s.db.QueryRowContext(context.Background(), selectProxyCountryIdByCode, c.CountryCode).Scan(&c.CountryId)
//}

//func (s *Store) InsertProxyCountry(c *pb.ProxyCountry) error {
//	err := s.db.QueryRowContext(context.Background(), insertProxyCountry, c.CountryName, c.CountryCode).Scan(&c.CountryId)
//	if err != nil {
//		return errors.Wrapf(err, "error insert proxy country %+v", c)
//	}
//	return nil
//}

//func (s *Store) GetNextProxyItemBatch(ctx context.Context, size int) error {
//	lastQueryTimeout := time.Now().UTC().Sub(s.nextProxyLastQuery) > nextProxyQueryTimeout
//	if !lastQueryTimeout {
//		s.log.Info("QueryTooOften")
//		return ErrQueryTooOften
//	}
//	s.nextProxyLastQuery = time.Now().UTC()
//	rows, err := s.db.QueryContext(ctx, getNextProxyItemBatch, size)
//	if err != nil {
//		return err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var ip net.IP
//		p := pb.ProxyItem{}
//		err := rows.Scan(&p.ProxyId, &ip, &p.ProxyPort)
//		if err != nil {
//			s.log.Error(errors.Wrap(err, "error wile scan getNextProxyItemBatch"))
//			continue
//		}
//		p.ProxyIp = ip.String()
//		s.nextProxyQueue <- p
//	}
//	return nil
//}

//func (s *Store) GetNextProxyItem(ctx context.Context) (*pb.ProxyItem, error) {
//	proxyItem, err := s.NextProxyItemProducer(ctx)
//	switch err {
//	case nil:
//		return &proxyItem, nil
//	case pgx.ErrNoRows:
//		return nil, ErrNoRows
//	default:
//		return nil, errors.Wrap(err, "Store: GetNextProxyItem error")
//	}
//}

// NextProxyItemProducer вынимает ProxyItem из nextProxyQueue
// Если элемента нет, вызывает функцию пополнения, и повторяет попытку взять элемент
//func (s *Store) NextProxyItemProducer(ctx context.Context) (pb.ProxyItem, error) {
//	for {
//		select {
//		case p := <-s.nextProxyQueue:
//			return p, nil
//		default:
//			err := s.GetNextProxyItemBatch(ctx, 100)
//			switch err {
//			case nil:
//			case ErrQueryTooOften:
//				time.Sleep(100 * time.Millisecond)
//			default:
//				return pb.ProxyItem{}, err
//			}
//		}
//	}
//}

//func (s *Store) CreateProxyStat(ctx context.Context, stat *pb.ProxyStat) error {
//	return s.db.QueryRowContext(ctx, insertProxyStat, stat.ProxyId, stat.ConnTime, stat.ConnStatus).Scan(&stat.CreatedAt)
//}

//func (s *Store) GetNextProxyItem(p *ProxyItem) error {
//	var ip net.IP
//	err := s.db.QueryRow(context.Background(), getNextProxyItem).Scan(&p.ProxyId, &ip, &p.ProxyPort)
//	switch err {
//	case nil:
//		p.ProxyIp = ip.String()
//		return nil
//	case pgx.ErrNoRows:
//		return ErrNoRows
//	default:
//		return errors.Wrap(err, "Store: GetNextProxyItem error")
//	}
//}

//func (s *Store) InsertProxyItem(p *pb.ProxyItem) error {
//	return s.db.QueryRowContext(context.Background(), insertProxyItem,
//		p.ProxyIp, p.ProxyPort,
//		p.ProxyCountry.CountryId).Scan(&p.ProxyId)
//}
//func (s *Store) SelectProxyIdByUI(ctx context.Context, p *pb.ProxyItem) error {
//	return s.db.QueryRowContext(ctx, selectProxyIdByUI, p.ProxyIp, p.ProxyPort).Scan(&p.ProxyId)
//}
