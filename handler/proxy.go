// @author xiangqian
// @date 2025/08/10 13:51
package handler

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
)

func proxy(prefix string) {
	http.HandleFunc(fmt.Sprintf("%s/proxy/", prefix), func(w http.ResponseWriter, r *http.Request) {
		// 目标请求地址
		var url = r.Header.Get("X-Url")
		log.Println(url)

		// 创建一个新的请求
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
	})
}
