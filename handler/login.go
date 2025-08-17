// @author xiangqian
// @date 2025/08/16 02:35
package handler

import (
	"fmt"
	"net/http"
	"strings"
)

// 处理用户请求
func uhandle(prefix, user, passwd string) {
	if user == "" || passwd == "" {
		return
	}

	login(prefix)
	login1(prefix, user, passwd)
	logout(prefix)
}

func login(prefix string) {
	http.HandleFunc(fmt.Sprintf("%s/login", prefix), func(w http.ResponseWriter, r *http.Request) {
		// 判断会话是否已过期
		if !expired(r) {
			// 重定向到首页
			http.Redirect(w, r, fmt.Sprintf("%s/", prefix), http.StatusFound)
			return
		}

		user, _ := gcookie(r, "user")
		errstr, _ := gcookie(r, "error")
		var data = map[string]any{
			"prefix": prefix,
			"user":   user,
			"error":  errstr,
		}
		execTmpl("login", data, w)
	})
}

func login1(prefix, user, passwd string) {
	http.HandleFunc(fmt.Sprintf("%s/login1", prefix), func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		u := strings.TrimSpace(r.Form.Get("user"))
		p := strings.TrimSpace(r.Form.Get("passwd"))
		if u == user && p == passwd {
			ssession(w)
			http.Redirect(w, r, fmt.Sprintf("%s/", prefix), http.StatusFound)
			return
		}

		scookie(w, "user", u, 2)
		scookie(w, "error", "用户名或密码错误", 2)
		http.Redirect(w, r, fmt.Sprintf("%s/login", prefix), http.StatusFound)
	})
}

func logout(prefix string) {
	http.HandleFunc(fmt.Sprintf("%s/logout", prefix), func(w http.ResponseWriter, r *http.Request) {
		// 判断会话是否已过期
		if !expired(r) {
			dsession(w, r)
		}
		// 重定向到登录页
		http.Redirect(w, r, fmt.Sprintf("%s/login", prefix), http.StatusFound)
	})
}
