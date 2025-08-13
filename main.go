// @author xiangqian
// @date 2025/08/09 19:43
package main

import (
	"apic/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 初始化日志
	err := InitLog()
	if err != nil {
		log.Fatalln(err)
	}

	// 加载配置文件
	config, err := LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n", config)

	// 处理请求
	handler.Handle(config.Prefix, config.Dir)

	// 启动服务
	port := config.Port
	log.Printf("Server starting on port %d ...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}
