package impl

import (
	"context"
	"fmt"
	"github.com/zginkgo/ginkgo_cmdb/apps"
	"github.com/zginkgo/ginkgo_cmdb/apps/book"
	"github.com/zginkgo/ginkgo_cmdb/conf"
	"testing"
)

var (
	ins book.ServiceServer
)

func TestQueryBook(t *testing.T) {
	ss, err := ins.QueryBook(context.Background(), book.NewQueryBookRequest())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ss.Items)
	t.Log(ss)
}

func TestDescribeBook(t *testing.T) {
	ss, err := ins.DescribeBook(context.Background(), book.NewDescribeBookRequest("ccml6veg26u7h4qj4f20"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ss)
}

func TestCreateBook(t *testing.T) {
	req := book.NewCreateBookRequest()
	req.CreateBy = "youmen"
	req.Name = "三体"
	req.Author = "刘慈溪"
	ss, err := ins.CreateBook(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ss)
}

func init() {
	if err := conf.LoadConfigFromEnv(); err != nil {
		panic(err)
	}

	if err := apps.InitAllApp(); err != nil {
		panic(err)
	}

	ins = apps.GetGrpcApp(book.AppName).(book.ServiceServer)
}
