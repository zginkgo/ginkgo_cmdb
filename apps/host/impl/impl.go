package impl

import (
	"database/sql"
	"github.com/zginkgo/ginkgo_cmdb/apps"
	"github.com/zginkgo/ginkgo_cmdb/apps/host"
	"github.com/zginkgo/ginkgo_cmdb/conf"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"
)

var (
	// Service 服务实例
	svr = &service{}
)

type service struct {
	db  *sql.DB
	log logger.Logger
	host.UnimplementedServiceServer
}

func (s *service) Name() string {
	return host.AppName
}

func (s *service) Config() error {
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}

	s.log = zap.L().Named(s.Name())
	s.db = db
	return nil
}

func (s *service) Registry(server *grpc.Server) {
	host.RegisterServiceServer(server, svr)
}

func init() {
	apps.RegistryGrpcApp(svr)
}
