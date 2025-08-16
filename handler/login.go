// @author xiangqian
// @date 2025/08/16 02:35
package handler

import (
	"fmt"
	"net/http"
	"strings"
)

func login(prefix string, w http.ResponseWriter, r *http.Request) {
	user, _ := GetCookie(r, "user")
	errstr, _ := GetCookie(r, "error")
	var data = map[string]any{
		"prefix": prefix,
		"user":   user,
		"error":  errstr,
	}
	execTmpl("login", data, w)
}

func login1(prefix, user, passwd string, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ruser := strings.TrimSpace(r.Form.Get("user"))
	rpasswd := strings.TrimSpace(r.Form.Get("passwd"))
	if ruser == user && rpasswd == passwd {
		SetSession(w, nil)
		http.Redirect(w, r, fmt.Sprintf("%s/", prefix), http.StatusFound)
		return
	}

	SetCookie(w, "user", ruser, 2)
	SetCookie(w, "error", "用户名或密码错误", 2)
	http.Redirect(w, r, fmt.Sprintf("%s/login", prefix), http.StatusFound)
}

func logout(prefix string, w http.ResponseWriter, r *http.Request) {
	DelSession(w, r)
	http.Redirect(w, r, fmt.Sprintf("%s/login", prefix), http.StatusFound)
}
