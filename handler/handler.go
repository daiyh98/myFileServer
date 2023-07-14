package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		//返回上传html页面
		data, err := os.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		} else {
			io.WriteString(w, string(data))
		}
	} else if request.Method == "POST" {
		//接收文件流及存储到本地目录
		file, head, err := request.FormFile("file")
		if err != nil {
			fmt.Printf("failed to get data, err: %s\n", err.Error())
			return
		}
		defer file.Close()
		//创建新文件
		newFile, err := os.Create("/tmp/" + head.Filename)
		if err != nil {
			fmt.Printf("failed to create file, err: %s\n", err.Error())
			return
		}
		defer newFile.Close()
		//将内存中文件流的内容拷贝到新的文件中
		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("failed to save data into file, err: %s\n", err.Error())
		}

		http.Redirect(w, request, "/file/upload/succeed", http.StatusFound)
	}
}

// UploadSuccessHandler 上传已完成
func UploadSuccessHandler(w http.ResponseWriter, request *http.Request) {
	io.WriteString(w, "Upload succeeded!")
}
