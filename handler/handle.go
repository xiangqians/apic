// @author xiangqian
// @date 2025/08/09 20:00
package handler

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

// 是否是开发环境
const dev = true

//go:embed image/* css/* js/* html/*
var embedfs embed.FS

func Handle(prefix, user, passwd string) error {
	// 处理静态资源请求
	err := shandle(prefix)
	if err != nil {
		return err
	}

	// 处理用户请求
	uhandle(prefix, user, passwd)

	// 处理其他请求
	http.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(fmt.Sprintf("%s", prefix), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var path = r.URL.Path

		var suser string
		if user != "" && passwd != "" {
			// 判断会话是否已过期
			if expired(r) {
				// 重定向到登录页
				http.Redirect(w, r, fmt.Sprintf("%s/login", prefix), http.StatusFound)
				return
			}
			suser = user
		}

		// 首页
		if path == "" || path == "/" {
			index(prefix, suser, w)
			return
		}

		// 文档
		swagger(prefix, suser, w, r)
	})))

	// 处理未匹配请求
	if prefix != "" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// 判断会话是否已过期
			if user != "" && passwd != "" && expired(r) {
				// 重定向到登录页
				http.Redirect(w, r, fmt.Sprintf("%s/login", prefix), http.StatusFound)
				return
			}

			// 错误页
			erro(prefix, errors.New(http.StatusText(http.StatusNotFound)), w)
		})
	}

	// 处理代理请求
	phandle(prefix)

	return nil
}

// 错误页
func erro(prefix string, err error, w http.ResponseWriter) {
	var data = map[string]any{
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
