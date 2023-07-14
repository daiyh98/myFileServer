package handler

import (
	"fmt"
	"io"
	"myFileServer/metaInfo"
	"myFileServer/utils"
	"net/http"
	"os"
	"time"
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

		fileMeta := metaInfo.FileMeta{
			FileName: head.Filename,
			//FileSha1:     utils.FileSha1(file.(*os.File)),
			FileLocation: "/tmp/" + head.Filename,
			UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
			//FileSize:     head.Size,
		}

		//创建新文件
		newFile, err := os.Create(fileMeta.FileLocation)
		if err != nil {
			fmt.Printf("failed to create file, err: %s\n", err.Error())
			return
		}
		defer newFile.Close()
		//将内存中文件流的内容拷贝到新的文件中
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("failed to save data into file, err: %s\n", err.Error())
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = utils.FileSha1(newFile)
		metaInfo.UpdateFileMetaMap(fileMeta)

		http.Redirect(w, request, "/file/upload/succeed", http.StatusFound)
	}
}

// UploadSuccessHandler 上传已完成
func UploadSuccessHandler(w http.ResponseWriter, request *http.Request) {
	io.WriteString(w, "Upload succeeded!")
}
