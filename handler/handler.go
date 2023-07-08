package handler

import "net/http"

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回上传html页面
	} else if r.Method == "POST" {
		//接收文件流及存储到本地目录
	}
}
