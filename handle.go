// @author xiangqian
// @date 2025/08/09 20:00
package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed swagger.js swagger.css index.html
var embedfs embed.FS
var funcMap = template.FuncMap{
	"hasSuffix": strings.HasSuffix,
}
var tmpl *template.Template

var names map[string]string

func Handle(prefix, dir string) error {
	var err error
	tmpl, err = template.New("").Funcs(funcMap).ParseFS(embedfs, "index.html")
	if err != nil {
		return err
	}

	// 创建文件服务器
	handler := http.FileServer(http.FS(embedfs))
	// 处理请求
	http.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			entries, err := os.ReadDir(dir)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var name string
			names = make(map[string]string)
			for _, entry := range entries {
				if !entry.IsDir() {
					var fname = entry.Name()
					var jname = "swagger.json"
					if fname == "swagger.yaml" {
						if _, ok := names[jname]; ok {
							delete(names, jname)
						}
						name = fname
						names[fname] = dir + "/" + fname
					} else if fname == jname {
						if name == "" {
							name = fname
							names[fname] = dir + "/" + fname
						}
					} else if strings.HasSuffix(fname, ".js") && fname != "swagger.js" {
						names[fname] = dir + "/" + fname
					}
				}
			}
			log.Printf("names: %+v\n", names)

			var data = map[string]interface{}{
				"prefix": prefix,
				"name":   name,
				"names":  names,
			}

			err = tmpl.ExecuteTemplate(w, "index.html", data)
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

		var name = r.URL.Path[1:]
		if path, ok := names[name]; ok {
			if r.Method == http.MethodGet {
				data, err := os.ReadFile(path)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				var contentType string
				if strings.HasSuffix(name, ".yaml") {
					contentType = "application/yaml"
				} else if strings.HasSuffix(name, ".json") {
					contentType = "application/json"
				} else if strings.HasSuffix(name, ".js") {
					contentType = "application/javascript"
				}
				w.Header().Set("Content-Type", contentType)
				w.Write(data)
				return

			} else if r.Method == http.MethodPost && (strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".json")) {
				http.NotFound(w, r)
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

	return nil
}
