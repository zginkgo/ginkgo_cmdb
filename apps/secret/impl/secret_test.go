package impl

import (
	"context"
	"fmt"
	"github.com/zginkgo/ginkgo_cmdb/apps"
	"github.com/zginkgo/ginkgo_cmdb/apps/secret"
	"github.com/zginkgo/ginkgo_cmdb/conf"
	"testing"
)

var (
	ins secret.ServiceServer
)

func TestQuerySecret(t *testing.T) {
	ss, err := ins.QuerySecret(context.Background(), secret.NewQuerySecretRequest())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ss.Items)
	t.Log(ss)
}

func TestDescribeSecret(t *testing.T) {
	ss, err := ins.DescribeSecret(context.Background(), secret.NewDescribeSecretRequest("ccrbkr6g26ud9jbnrik0"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ss)
}

func TestCreateSecret(t *testing.T) {
	req := secret.NewCreateSecretRequest()
	req.Description = "测试用例"
	req.ApiKey = "213421432134"
	req.ApiSecret = "234134"
	req.AllowRegions = []string{"*"}
	ss, err := ins.CreateSecret(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ss)
}

func init() {
	if err := conf.LoadConfigFromEnv(); err != nil {
		panic(err)
	}

	// 全局日志对象初始化
	if err := apps.InitAllApp(); err != nil {
		panic(err)
	}

	ins = apps.GetGrpcApp(secret.AppName).(secret.ServiceServer)
}
