package protocol

import (
	"context"
	"github.com/zginkgo/ginkgo_cmdb/apps"
	"github.com/infraboard/mcenter/apps/instance"
	"github.com/infraboard/mcenter/client/rpc"
	"github.com/infraboard/mcenter/client/rpc/lifecycle"
	"google.golang.org/grpc"
	"net"
	"time"

	"github.com/infraboard/mcube/grpc/middleware/recovery"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/zginkgo/ginkgo_cmdb/conf"
)

// NewGRPCService todo
func NewGRPCService() *GRPCService {
	log := zap.L().Named("GRPC Service")

	rc := recovery.NewInterceptor(recovery.NewZapRecoveryHandler())
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		rc.UnaryServerInterceptor(),
	))

	// 控制Grpc 启动其他服务, 比如注册中心
	ctx, cancel := context.WithCancel(context.Background())

	return &GRPCService{
		svr:    grpcServer,
		l:      log,
		c:      conf.C(),
		ctx:    ctx,
		cancel: cancel,
	}
}

// GRPCService grpc服务
type GRPCService struct {
	svr    *grpc.Server
	l      logger.Logger
	c      *conf.Config
	ctx    context.Context
	cancel context.CancelFunc
	lf     lifecycle.Lifecycler
}

// Start 启动GRPC服务
func (s *GRPCService) Start() {
	// 装载所有GRPC服务
	apps.LoadGrpcApp(s.svr)

	// 启动HTTP服务
	lis, err := net.Listen("tcp", s.c.App.GRPC.Addr())
	if err != nil {
		s.l.Errorf("listen grpc tcp conn error, %s", err)
		return
	}

	time.AfterFunc(1*time.Second, s.registry)
	s.l.Infof("GRPC 服务监听地址: %s", s.c.App.GRPC.Addr())
	if err := s.svr.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			s.l.Info("service is stopped")
		}

		s.l.Error("start grpc service error, %s", err.Error())
		return
	}
}

func (s *GRPCService) registry() {
	req := instance.NewRegistryRequest()
	req.Address = s.c.App.GRPC.Addr()
	lf, err := rpc.C().Registry(s.ctx, req)
	if err != nil {
		s.l.Errorf("registry to mcenter error, %s", err)
		return
	}
	s.l.Infof("registry to mcenter success")

	// 注销时候需要使用
	s.lf = lf
}

// Stop 启动GRPC服务
func (s *GRPCService) Stop() error {
	if s.lf != nil {
		if err := s.lf.UnRegistry(s.ctx); err != nil {
			s.l.Errorf("unregistry error, %s", err)
		} else {
			s.l.Infof("unregistry success")
		}
	}
	s.svr.GracefulStop()
	return nil
}
