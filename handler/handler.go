package handler

import (
	"io"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回上传html页面
		data, err := os.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		} else {
			io.WriteString(w, string(data))
		}
	} else if r.Method == "POST" {
		//接收文件流及存储到本地目录
	}
}
