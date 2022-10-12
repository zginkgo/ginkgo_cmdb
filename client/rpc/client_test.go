package rpc_test

import (
	"context"
	"fmt"
	mcenter "github.com/infraboard/mcenter/client/rpc"
	"github.com/zginkgo/ginkgo_cmdb/apps/secret"
	"github.com/zginkgo/ginkgo_cmdb/client/rpc"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	should := assert.New(t)
	conf := mcenter.NewDefaultConfig()

	// 传递Mcenter配置, 客户端通过Mcenter进行搜索
	c, err := rpc.NewClient(conf)

	if should.NoError(err) {
		//rs, err := c.Resource().Search(context.Background(), resource.NewSearchRequest())
		rs, err := c.Secret().QuerySecret(context.Background(), secret.NewQuerySecretRequest())
		should.NoError(err)
		fmt.Println(rs)
	}
}

func init() {
	// 提前加载好 mcenter 客户端, resolver 需要使用
	err := mcenter.LoadClientFromEnv()
	if err != nil {
		panic(err)
	}
}
