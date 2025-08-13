// @author xiangqian
// @date 2025/08/10 11:34
package handler

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func d(name, dir string, w http.ResponseWriter, r *http.Request) {
	names, err := readDir(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tname string
	for _, nam := range names {
		if nam == name {
			tname = nam
			break
		}
	}
	if tname == "" {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		get(tname, dir, w, r)
		return
	}

	if r.Method == http.MethodPost {
		post(tname, dir, w, r)
		return
	}

	http.NotFound(w, r)
}

func get(name, dir string, w http.ResponseWriter, r *http.Request) {
	var contentType string
	if strings.HasSuffix(name, ".yaml") {
		contentType = "application/yaml"
	} else if strings.HasSuffix(name, ".json") {
		contentType = "application/json"
	} else if strings.HasSuffix(name, ".js") {
		contentType = "application/javascript"
	}
	if contentType == "" {
		http.NotFound(w, r)
		return
	}

	data, err := os.ReadFile(path.Join(dir, name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

func post(name, dir string, w http.ResponseWriter, r *http.Request) {
	if !(strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".json")) {
		http.NotFound(w, r)
		return
	}

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

	file, err := os.Create(path.Join(dir, name))
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
	w.Write([]byte("ok"))
}

func readDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var names = make([]string, 0, len(entries))
	names = append(names, "")
	for _, entry := range entries {
		var name = entry.Name()
		if strings.HasSuffix(name, ".yaml") {
			names[0] = name
		} else if strings.HasSuffix(name, ".json") {
			if names[0] == "" {
				names[0] = name
			}
		} else if strings.HasSuffix(name, ".js") {
			names = append(names, name)
		}
	}
	return names, nil
}
