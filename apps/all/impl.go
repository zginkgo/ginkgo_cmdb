package all

import (
	_ "github.com/zginkgo/ginkgo_cmdb/apps/book/impl"
	_ "github.com/zginkgo/ginkgo_cmdb/apps/host/impl"
	_ "github.com/zginkgo/ginkgo_cmdb/apps/resource/impl"
	// 注册所有GRPC服务模块, 暴露给框架GRPC服务器加载, 注意 导入有先后顺序
	_ "github.com/zginkgo/ginkgo_cmdb/apps/secret/impl"
	_ "github.com/zginkgo/ginkgo_cmdb/apps/task/impl"
)
