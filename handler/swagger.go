// @author xiangqian
// @date 2025/08/10 11:34
package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// http.Error(

func swagger(prefix string, w http.ResponseWriter, r *http.Request) {
	var view bool
	var path = r.URL.Path
	if strings.HasSuffix(path, "/view") {
		path = path[0 : len(path)-len("/view")]
		view = true
	}
	arr := strings.Split(path, "/")
	if len(arr) != 3 {
		erro(prefix, errors.New(http.StatusText(http.StatusNotFound)), w)
		return
	}

	dirs, err := readDirs()
	if err != nil {
		erro(prefix, err, w)
		return
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
		erro(prefix, errors.New(http.StatusText(http.StatusNotFound)), w)
		return
	}

	var file File
	var fid = arr[2]
	if fid != "" {
		for _, f := range dir.Files {
			if f.Id == fid {
				file = f
				break
			}
		}
		if file.Id == "" {
			erro(prefix, errors.New(http.StatusText(http.StatusNotFound)), w)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		if view {
			vswagger(prefix, dir, file, w)
		} else {
			gswagger(prefix, dir, file, w)
		}
	case http.MethodPost:
		pswagger(prefix, dir, file, w, r)
	default:
		erro(prefix, errors.New(http.StatusText(http.StatusNotFound)), w)
	}
}

func vswagger(prefix string, dir Dir, file File, w http.ResponseWriter) {
	var data = map[string]interface{}{
		"prefix": prefix,
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

func pswagger(prefix string, dir Dir, file File, w http.ResponseWriter, r *http.Request) {

}

//func fdir(edir, efile string, w http.ResponseWriter, r *http.Request) {
//	ddir, err := decode(edir)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	dirs, err := readDirs()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	var tdir string
//	for _, dir := range dirs {
//		if dir == ddir {
//			tdir = dir
//			break
//		}
//	}
//	if tdir == "" {
//		http.Error(w, "empty", http.StatusInternalServerError)
//		return
//	}
//
//	if efile != "" {
//		dfile, err := decode(efile)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		files, err := readFiles(tdir)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		var tfile string
//		for _, file := range files {
//			if file == dfile {
//				tfile = file
//				break
//			}
//		}
//		if tfile == "" {
//			http.NotFound(w, r)
//			return
//		}
//
//		//if r.Method == http.MethodGet {
//		//	get(tname, dir, w, r)
//		//	return
//		//}
//		//
//		//if r.Method == http.MethodPost {
//		//	post(tname, dir, w, r)
//		//	return
//		//}
//	}
//
//	http.NotFound(w, r)
//}
//
//func post(name, dir string, w http.ResponseWriter, r *http.Request) {
//	if !(strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".json")) {
//		http.NotFound(w, r)
//		return
//	}
//
//	if r.Header.Get("Content-Type") != "text/plain" {
//		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
//		return
//	}
//
//	data, err := io.ReadAll(r.Body)
//	defer r.Body.Close()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	file, err := os.Create(path.Join(dir, name))
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer file.Close()
//
//	writer := bufio.NewWriterSize(file, 1<<12) // 4KB 缓冲区
//	_, err = writer.Write(data)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	err = writer.Flush()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "text/plain")
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte("ok"))
//}
