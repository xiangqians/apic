// @author xiangqian
// @date 2025/08/09 20:00
package handler

import (
	"embed"
	"fmt"
	"net/http"
	"strings"
)

//go:embed v2/* v3/*
var embedfs embed.FS

func Handle(prefix, dir string) {
	// 文件服务器
	handler := http.FileServer(http.FS(embedfs))

	// 处理请求
	http.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(fmt.Sprintf("%s", prefix), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var path = r.URL.Path

		// index
		if path == "" || path == "/" || path == "/edit" || path == "/edit/" {
			index(prefix, dir, embedfs, w, r)
			return
		}

		// file
		if strings.HasPrefix(path, "/v2/") || strings.HasPrefix(path, "/v3/") {
			handler.ServeHTTP(w, r)
			return
		}

		// dir
		d(path[1:], dir, w, r)
	})))

	// 代理
	proxy(prefix)
}
