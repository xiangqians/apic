// @author xiangqian
// @date 2025/08/10 11:34
package handler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func swagger(prefix, user string, w http.ResponseWriter, r *http.Request) {
	dir, file, t, err := parse(r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		if err != nil {
			erro(prefix, errors.New(http.StatusText(http.StatusNotFound)), w)
			return
		}
		if t {
			tswagger(prefix, user, dir, file, w)
		} else {
			gswagger(prefix, dir, file, w)
		}

	case http.MethodPost:
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(err.Error()))
			return
		}
		pswagger(dir, file, w, r)

	default:
		erro(prefix, errors.New(http.StatusText(http.StatusNotFound)), w)
	}
}

func parse(path string) (Dir, File, bool, error) {
	var t bool
	if strings.HasSuffix(path, "/t") {
		path = path[0 : len(path)-len("/t")]
		t = true
	}

	arr := strings.Split(path, "/")
	if len(arr) != 3 {
		return Dir{}, File{}, false, errors.New("!=3")
	}

	dirs, err := readDirs()
	if err != nil {
		return Dir{}, File{}, false, err
	}

	var dir Dir
	var did = arr[1]
	for _, d := range dirs {
		if d.Id == did {
			dir = d
			break
		}
	}
	if dir.Id == "" {
		return Dir{}, File{}, false, errors.New("dir not found")
	}

	var file File
	var fid = arr[2]
	for _, f := range dir.Files {
		if f.Id == fid {
			file = f
			break
		}
	}
	if file.Id == "" {
		return Dir{}, File{}, false, errors.New("file not found")
	}

	return dir, file, t, nil
}

func tswagger(prefix, user string, dir Dir, file File, w http.ResponseWriter) {
	var data = map[string]any{
		"prefix": prefix,
		"user":   user,
		"dir":    dir,
		"file":   file,
	}
	execTmpl("swagger", data, w)
}

func gswagger(prefix string, dir Dir, file File, w http.ResponseWriter) {
	var contentType string
	switch file.Type {
	case TypeYaml:
		contentType = "application/yaml"
	case TypeJson:
		contentType = "application/json"
	case TypeJs:
		contentType = "application/javascript"
	}
	if contentType == "" {
		erro(prefix, errors.New(http.StatusText(http.StatusNotFound)), w)
		return
	}

	data, err := os.ReadFile(path.Join(dir.Name, file.Name))
	if err != nil {
		erro(prefix, err, w)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", strings.ReplaceAll(url.QueryEscape(file.Name), "+", "%20")))
	w.Write(data)
}

func pswagger(dir Dir, file File, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	if file.Type != TypeYaml && file.Type != TypeJson {
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}

	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	f, err := os.Create(path.Join(dir.Name, file.Name))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer f.Close()

	writer := bufio.NewWriterSize(f, 1<<12) // 4KB 缓冲区
	_, err = writer.Write(data)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	err = writer.Flush()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("ok"))
}
