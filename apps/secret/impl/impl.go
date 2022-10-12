package impl

import (
	"database/sql"
	"github.com/zginkgo/ginkgo_cmdb/apps"
	"github.com/zginkgo/ginkgo_cmdb/apps/secret"
	"github.com/zginkgo/ginkgo_cmdb/conf"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"
)

var (
	svr = &service{}
)

type service struct {
	db     *sql.DB
	log    logger.Logger
	secret secret.ServiceServer
	secret.UnimplementedServiceServer
}

func (s *service) Name() string {
	return secret.AppName
}

func (s *service) Config() error {
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}

	s.log = zap.L().Named(s.Name())
	s.db = db
	s.secret = apps.GetGrpcApp(secret.AppName).(secret.ServiceServer)
	return nil
}

func (s *service) Registry(server *grpc.Server) {
	secret.RegisterServiceServer(server, svr)
}

func init() {
	apps.RegistryGrpcApp(svr)
}
