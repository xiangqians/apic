// @author xiangqian
// @date 2025/08/09 20:00
package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

//go:embed swagger.js swagger.css index.html
var embedfs embed.FS

func Handle(prefix, name string) {
	// 创建文件服务器
	handler := http.FileServer(http.FS(embedfs))
	// 处理请求
	http.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			var tname = "index.html"
			tmpl, err := template.New("").ParseFS(embedfs, tname)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var data = map[string]interface{}{
				"prefix": prefix,
				"name":   name,
			}
			err = tmpl.ExecuteTemplate(w, tname, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}

		if r.URL.Path == "/index.html" {
			http.Redirect(w, r, prefix+"/", http.StatusMovedPermanently)
			return
		}

		if r.URL.Path == "/"+name {
			if r.Method == http.MethodGet {
				data, err := os.ReadFile(name)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(data)
				return

			} else if r.Method == http.MethodPost {

				return
			}

			http.NotFound(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	})))

	// 处理代理请求
	http.Handle(fmt.Sprintf("%s/proxy", prefix), http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 目标服务器地址
		url := "http://example.com" + r.URL.Path
		req, err := http.NewRequest(r.Method, url, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		// 复制请求头
		for name, values := range r.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// 复制响应头
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		// 设置状态码
		w.WriteHeader(resp.StatusCode)

		// 复制响应体
		io.Copy(w, resp.Body)
	})))
}
