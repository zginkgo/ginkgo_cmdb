package impl

import (
	"database/sql"
	"github.com/zginkgo/ginkgo_cmdb/apps"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"

	"github.com/zginkgo/ginkgo_cmdb/apps/book"
	"github.com/zginkgo/ginkgo_cmdb/conf"
)

var (
	// Service 服务实例
	svr = &service{}
)

type service struct {
	db *sql.DB

	log  logger.Logger
	book book.ServiceServer
	book.UnimplementedServiceServer
}

func (s *service) Config() error {
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}
	s.db = db

	s.log = zap.L().Named(s.Name())
	s.book = apps.GetGrpcApp(book.AppName).(book.ServiceServer)
	return nil
}

func (s *service) Name() string {
	return book.AppName
}

func (s *service) Registry(server *grpc.Server) {
	book.RegisterServiceServer(server, svr)
}

func init() {
	apps.RegistryGrpcApp(svr)
}
