// @author xiangqian
// @date 2025/08/09 19:55
package main

import pkg_ini "gopkg.in/ini.v1"

// LoadConfig 加载配置文件
func LoadConfig() (Config, error) {
	var config Config
	file, err := pkg_ini.Load("config.ini")
	if err != nil {
		return config, err
	}

	section, err := file.GetSection("")
	if err != nil {
		return config, err
	}
	config.Port = uint16(section.Key("port").MustUint())
	config.Prefix = section.Key("prefix").String()
	config.Name = section.Key("name").String()

	return config, nil
}

type Config struct {
	Port   uint16 // HTTP 监听端口
	Prefix string // HTTP 请求前缀
	Name   string // Swagger 文件地址
}
