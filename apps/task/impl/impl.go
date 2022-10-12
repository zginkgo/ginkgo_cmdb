package impl

import (
	"database/sql"
	"github.com/zginkgo/ginkgo_cmdb/apps"
	"github.com/zginkgo/ginkgo_cmdb/apps/host"
	"github.com/zginkgo/ginkgo_cmdb/apps/secret"
	"github.com/zginkgo/ginkgo_cmdb/apps/task"
	"github.com/zginkgo/ginkgo_cmdb/conf"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"
)

var (
	svr = &impl{}
)

type impl struct {
	db  *sql.DB
	log logger.Logger
	task.UnimplementedServiceServer

	secret secret.ServiceServer
	host   host.ServiceServer
}

func (s *impl) Config() error {
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}

	s.log = zap.L().Named(s.Name())
	s.db = db
	//s.secret = apps.GetGrpcApp(secret.AppName).(secret.ServiceServer)

	// 通过mock 来解耦 s.secret = &secretMock{}
	s.secret = &secretMock{}
	s.secret = apps.GetGrpcApp(secret.AppName).(secret.ServiceServer)
	s.host = apps.GetGrpcApp(host.AppName).(host.ServiceServer)
	return nil
}

func (s *impl) Name() string {
	return task.AppName
}

func (s *impl) Registry(server *grpc.Server) {
	task.RegisterServiceServer(server, svr)
}

func init() {
	apps.RegistryGrpcApp(svr)
}
