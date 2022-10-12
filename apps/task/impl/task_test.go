package impl

import (
	"context"
	"github.com/zginkgo/ginkgo_cmdb/apps"
	_ "github.com/zginkgo/ginkgo_cmdb/apps/all"
	"github.com/zginkgo/ginkgo_cmdb/apps/resource"
	"github.com/zginkgo/ginkgo_cmdb/apps/task"
	"github.com/zginkgo/ginkgo_cmdb/conf"
	"github.com/infraboard/mcube/logger/zap"
	"testing"
)

var (
	ins  task.ServiceServer
	mock secretMock
)

//func TestDescribeSecret(t *testing.T) {
//	req := secret.NewDescribeSecretRequest("1")
//	mock.DescribeSecret(context.Background(), req)
//	fmt.Println(1234)
//}

func TestCreateTask(t *testing.T) {
	//secretReq := secret.NewDescribeSecretRequest("1")
	//res, _ := mock.DescribeSecret(context.Background(), secretReq)
	//fmt.Println(res.Data, "---->")
	req := task.NewCreateTaskRequest()
	req.Type = task.Type_RESOURCE_SYNC
	req.Region = "ap-nanjing"
	req.ResourceType = resource.Type_HOST
	req.SecretId = "cd0ddf6g26u2e7b7khs0"
	taskIns, err := ins.CreateTask(context.Background(), req)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(taskIns)
}

func init() {
	// 通过环境变量加载测试配置
	if err := conf.LoadConfigFromEnv(); err != nil {
		panic(err)
	}

	// 全局日志对象初始化
	zap.DevelopmentSetup()

	// 初始化所有实例
	if err := apps.InitAllApp(); err != nil {
		panic(err)
	}

	ins = apps.GetGrpcApp(task.AppName).(task.ServiceServer)
}
