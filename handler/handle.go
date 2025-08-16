// @author xiangqian
// @date 2025/08/09 20:00
package handler

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
)

// 是否是开发环境
const dev = true

//go:embed css/* js/* html/*
var embedfs embed.FS

func Handle(prefix string) error {
	// 文件服务器
	for _, dir := range []string{"image", "css", "js"} {
		var iofs fs.FS
		if dev {
			// 从文件系统加载，支持热重载
			iofs = os.DirFS(fmt.Sprintf("handler/%s", dir))
		} else {
			// 使用 fs.Sub 获取子文件系统
			var err error
			iofs, err = fs.Sub(embedfs, dir)
			if err != nil {
				return err
			}
		}

		handler := http.FileServer(http.FS(iofs))
		var pattern = fmt.Sprintf("%s/%s/", prefix, dir)
		http.Handle(pattern, http.StripPrefix(pattern, handler))
	}

	// 处理请求
	http.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(fmt.Sprintf("%s", prefix), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var path = r.URL.Path
		if path == "" || path == "/" {
			index(prefix, w)
			return
		}

		swagger(prefix, w, r)
	})))

	if prefix != "" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// 未匹配路由返回错误页
			erro(prefix, errors.New(http.StatusText(http.StatusNotFound)), w)
		})
	}

	// 代理
	proxy(prefix)

	return nil
}

// 错误页
func erro(prefix string, err error, w http.ResponseWriter) {
	var data = map[string]interface{}{
		"prefix": prefix,
		"error":  err,
	}
	execTmpl("error", data, w)
}

// 执行模板
func execTmpl(name string, data any, w http.ResponseWriter) {
	var err error
	var tmpl *template.Template

	// 解析模板
	if dev {
		// 从文件系统加载，支持热重载
		tmpl, err = template.New("").ParseGlob("handler/html/*.html")
	} else {
		tmpl, err = template.New("").ParseFS(embedfs, "html/*.html")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 执行模板
	err = tmpl.ExecuteTemplate(w, fmt.Sprintf("%s.html", name), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
