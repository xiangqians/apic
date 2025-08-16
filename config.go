// @author xiangqian
// @date 2025/08/09 19:55
package main

import (
	"gopkg.in/ini.v1"
	"strings"
)

// LoadConfig 加载配置文件
func LoadConfig() (Config, error) {
	var config Config
	file, err := ini.Load("config.ini")
	if err != nil {
		return config, err
	}

	section, err := file.GetSection("")
	if err != nil {
		return config, err
	}

	config.Port = uint16(section.Key("port").MustUint())
	config.Prefix = strings.TrimSpace(section.Key("prefix").String())
	config.User = strings.TrimSpace(section.Key("user").String())
	config.Passwd = strings.TrimSpace(section.Key("passwd").String())
	return config, nil
}

type Config struct {
	Port   uint16 // HTTP 监听端口
	Prefix string // HTTP 请求前缀
	User   string // 登录用户
	Passwd string // 登录密码
}
