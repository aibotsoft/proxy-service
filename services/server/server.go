package server

import (
	"context"
	"database/sql"
	pb "github.com/aibotsoft/gen/proxypb"
	"github.com/aibotsoft/micro/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"strconv"
)

type Server struct {
	cfg   *config.Config
	log   *zap.SugaredLogger
	store *Store
	gs    *grpc.Server
	//pb.UnimplementedProxyServer
}

func NewServer(cfg *config.Config, log *zap.SugaredLogger, db *sql.DB) *Server {
	return &Server{
		cfg:   cfg,
		log:   log,
		store: NewStore(cfg, log, db),
		gs:    grpc.NewServer(),
	}
}

func (s *Server) GracefulStop() {
	s.log.Debug("begin proxy server gracefulStop")
	s.gs.GracefulStop()
	s.log.Debug("end proxy server gracefulStop")
}

func (s *Server) Serve() error {
	addr := net.JoinHostPort("", strconv.Itoa(s.cfg.ProxyService.GrpcPort))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "net.Listen error")
	}
	pb.RegisterProxyServer(s.gs, s)
	s.log.Info("gRPC Proxy Server listens on port ", strconv.Itoa(s.cfg.ProxyService.GrpcPort))
	return s.gs.Serve(lis)
}

func (s *Server) CreateProxyStat(ctx context.Context, req *pb.CreateProxyStatRequest) (*pb.CreateProxyStatResponse, error) {
	proxyStat := req.GetProxyStat()
	err := s.store.CreateProxyStat(ctx, proxyStat)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateProxyStat: %v", err)
	}
	return &pb.CreateProxyStatResponse{ProxyStat: proxyStat}, nil
}

func (s *Server) CreateProxy(ctx context.Context, req *pb.CreateProxyRequest) (*pb.CreateProxyResponse, error) {
	proxyItem := req.GetProxyItem()
	err := s.store.GetOrCreateProxyItem(ctx, proxyItem)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetOrCreateProxyItem: %v", err)
	}
	return &pb.CreateProxyResponse{ProxyItem: proxyItem}, nil
}

// GetNextProxy возвращает прокси которое нужно проверить.
// Возвращаются те которые еще не проверялись, либо отсортированные по времени проверки.
func (s *Server) GetNextProxy(ctx context.Context, req *pb.GetNextProxyRequest) (*pb.GetNextProxyResponse, error) {
	proxyItem, err := s.store.GetNextProxyItem(ctx)
	switch err {
	case nil:
		return &pb.GetNextProxyResponse{ProxyItem: proxyItem}, nil
	case ErrNoRows:
		return nil, status.Error(codes.NotFound, "no new proxy for check to return")
	default:
		return nil, status.Errorf(codes.Internal, "GetNextProxy error %v", err)
	}
}
func (s *Server) GetBestProxy(ctx context.Context, req *pb.GetBestProxyRequest) (*pb.GetBestProxyResponse, error) {
	proxyItem, err := s.store.GetBestProxy(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "store.GetBestProxy error: %v", err)
	}
	return &pb.GetBestProxyResponse{ProxyItem: proxyItem}, nil
}
