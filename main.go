package main

import (
	_ "github.com/curltech/go-colla-core/cache"
	_ "github.com/curltech/go-colla-stock/poem/controller"
	_ "github.com/curltech/go-colla-stock/stock/controller"
	"github.com/curltech/go-colla-web/app"
	_ "github.com/curltech/go-colla-web/basecode/controller"
	/**
	  引入包定义，执行对应包的init函数，从而引入某功能，在init函数根据初始化参数配置觉得是否启动该功能
	*/
	_ "github.com/curltech/go-colla-core/content"
	_ "github.com/curltech/go-colla-core/repository/search"
)

func main() {
	app.Start()
}
