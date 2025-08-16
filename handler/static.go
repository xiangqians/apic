// @author xiangqian
// @date 2025/08/15 21:22
package handler

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
)

// 处理静态资源请求
func shandle(prefix string) error {
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

	return nil
}
