// @author xiangqian
// @date 2025/08/09 20:00
package main

import (
	"bufio"
	"crypto/tls"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//go:embed swagger.js swagger.css index.html
var embedfs embed.FS

var tmpl *template.Template

var names map[string]string

// Handle 处理 HTTP 请求
func Handle(prefix, dir string) error {
	var err error
	tmpl, err = template.New("").Funcs(template.FuncMap{"hasSuffix": strings.HasSuffix}).ParseFS(embedfs, "index.html")
	if err != nil {
		return err
	}

	// 创建文件服务器
	fhandler := http.FileServer(http.FS(embedfs))
	// 处理 HTTP 请求
	http.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/edit" {
			index(prefix, dir, w, r)
			return
		}

		if r.URL.Path == "/index.html" {
			http.Redirect(w, r, fmt.Sprintf("%s/", prefix), http.StatusMovedPermanently)
			return
		}

		var name = r.URL.Path[1:]
		if path, ok := names[name]; ok {
			spec(name, path, w, r)
			return
		}

		fhandler.ServeHTTP(w, r)
	})))

	proxy(prefix)

	example(prefix)

	return nil
}

func index(prefix, dir string, w http.ResponseWriter, r *http.Request) {
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
		"edit":   r.URL.Path == "/edit",
		"prefix": prefix,
		"name":   name,
		"names":  names,
	}

	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Spec 目录
func spec(name, path string, w http.ResponseWriter, r *http.Request) {
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
	}

	if r.Method == http.MethodPost && (name == "swagger.yaml" || name == "swagger.json") {
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}

		data, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		file, err := os.Create(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		writer := bufio.NewWriterSize(file, 1<<12) // 4KB 缓冲区
		_, err = writer.Write(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = writer.Flush()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
		return
	}

	http.NotFound(w, r)
}

// 处理代理请求
func proxy(prefix string) {
	http.Handle(fmt.Sprintf("%s/proxy/", prefix), http.StripPrefix(fmt.Sprintf("%s/proxy", prefix), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var url = r.Header.Get("X-Url")
		log.Println(url)

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
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 跳过验证
		}
		client := &http.Client{Transport: transport}
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

func example(prefix string) {
	http.HandleFunc(fmt.Sprintf("%s/example/", prefix), func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(map[string]any{
			"code": "ok",
			"msg":  "Ok",
			"data": map[string]any{
				"path": r.URL.Path,
				"time": time.Now(),
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	})
}
