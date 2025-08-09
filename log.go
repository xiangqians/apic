// @author xiangqian
// @date 2025/08/09 19:56
package main

import (
	"log"
	"os"
	"time"
)

// InitLogger 初始化日志记录器
func InitLogger() error {
	// 将日志输出重定向到 stdout
	//log.SetOutput(os.Stdout)

	// 禁用默认的日期和时间前缀
	log.SetFlags(0)
	// 自定义日志输出
	log.SetOutput(&Logger{})

	return nil
}

type Logger struct{}

func (logger *Logger) Write(p []byte) (n int, err error) {
	return os.Stdout.WriteString(time.Now().Format("2006/01/02 15:04:05.000") + " " + string(p))
}
