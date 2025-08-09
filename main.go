// @author xiangqian
// @date 2025/08/09 19:43
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 初始化日志记录器
	InitLogger()

	// 加载配置文件
	config, err := LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n", config)

	// HTTP 请求处理
	Handle(config.Prefix, config.Dir)

	// 启动 HTTP 服务
	port := config.Port
	log.Printf("Server starting on port %d ...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}
