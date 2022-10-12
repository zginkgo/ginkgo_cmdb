package rpc

import (
	"fmt"
	"github.com/zginkgo/ginkgo_cmdb/apps/secret"
	"github.com/infraboard/mcenter/client/rpc"
	"github.com/infraboard/mcenter/client/rpc/auth"
	"github.com/infraboard/mcenter/client/rpc/resolver"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client *ClientSet
)

func SetGlobal(cli *ClientSet) {
	client = cli
}

// C Global
func C() *ClientSet {
	return client
}

// NewClient 传递注册中心地址
func NewClient(conf *rpc.Config) (*ClientSet, error) {
	zap.DevelopmentSetup()
	log := zap.L()

	conn, err := grpc.Dial(
		fmt.Sprintf("%s://%s", resolver.Scheme, "cmdb"), // Dial to "mcenter://cmdb"
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(auth.NewAuthentication(conf.ClientID, conf.ClientSecret)),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &ClientSet{
		conn: conn,
		log:  log,
	}, nil
}

// ClientSet  客户端
type ClientSet struct {
	conn *grpc.ClientConn
	log  logger.Logger
}

func (c *ClientSet) Secret() secret.ServiceClient {
	return secret.NewServiceClient(c.conn)
}
