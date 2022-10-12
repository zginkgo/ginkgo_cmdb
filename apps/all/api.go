package all

import (
	_ "github.com/zginkgo/ginkgo_cmdb/apps/book/api"
	_ "github.com/zginkgo/ginkgo_cmdb/apps/host/api"
	_ "github.com/zginkgo/ginkgo_cmdb/apps/resource/api"
	// 注册所有HTTP服务模块, 暴露给框架HTTP服务器加载
	_ "github.com/zginkgo/ginkgo_cmdb/apps/secret/api"
	_ "github.com/zginkgo/ginkgo_cmdb/apps/task/api"
)
