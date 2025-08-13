// @author xiangqian
// @date 2025/08/10 23:17
package handler

import (
	"bufio"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"regexp"
)

func index(prefix, dir string, embedfs embed.FS, w http.ResponseWriter, r *http.Request) {
	names, err := readDir(dir)
	if err != nil || len(names) == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := names[0]
	file, err := os.Open(path.Join(dir, name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var ver string
	re := regexp.MustCompile(`(?:"swagger"|swagger)\s*:\s*"([0-9.]+)"`)

	// Scanner 默认初始缓冲区大小 4KB，最大行长度限制 64KB
	scanner := bufio.NewScanner(file)
	// 可能因错误返回 false
	for scanner.Scan() {
		match := re.FindStringSubmatch(scanner.Text())
		if len(match) > 1 {
			ver = match[1]
			break
		}
	}
	// 循环结束后检查是否因错误退出
	if err = scanner.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ver == "2.0" {
		ver = "v2"
	} else if ver == "3.0" {
		ver = "v3"
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// HTML 模板
	var tname = fmt.Sprintf("%s/index.html", ver)
	tmpl, err := template.New("").ParseFS(embedfs, tname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data = map[string]interface{}{
		"edit":   r.URL.Path == "/edit",
		"prefix": prefix,
		"api":    name,
		"jss":    names[1:],
	}
	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
