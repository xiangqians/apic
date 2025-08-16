// @author xiangqian
// @date 2025/08/10 23:17
package handler

import (
	"net/http"
)

func index(prefix, user string, w http.ResponseWriter) {
	dirs, err := readDirs()
	if err != nil {
		erro(prefix, err, w)
		return
	}

	var data = map[string]any{
		"prefix": prefix,
		"user":   user,
		"dirs":   dirs,
	}
	execTmpl("index", data, w)
}
