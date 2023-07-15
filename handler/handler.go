package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"myFileServer/metaInfo"
	"myFileServer/utils"
	"net/http"
	"os"
	"time"
)

// UploadHandler 上传文件
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

// GetFileMetaHandler 查询文件元信息
func GetFileMetaHandler(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	fileHash := request.Form["filehash"][0]
	fileMeta := metaInfo.GetFileMeta(fileHash)
	data, err := json.Marshal(fileMeta)
	if err != nil {
		//fmt.Printf("failed to get file meta, err: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler 下载文件
func DownloadHandler(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	fileHash := request.Form["filehash"][0]
	fileMeta := metaInfo.GetFileMeta(fileHash)

	file, err := os.Open(fileMeta.FileLocation) //得到文件勾柄
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
	w.Write(data)
}

// UpdateFileMetaHandler 更新元信息（重命名）
func UpdateFileMetaHandler(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	opType := request.Form.Get("op")
	fileHash := request.Form.Get("filehash")
	newFileName := request.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if request.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := metaInfo.GetFileMeta(fileHash)
	curFileMeta.FileName = newFileName
	metaInfo.UpdateFileMetaMap(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// DeleteHandler 删除文件
func DeleteHandler(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fileSha1 := request.Form.Get("filehash")

	fileMeta := metaInfo.GetFileMeta(fileSha1)
	os.Remove(fileMeta.FileLocation)
	
	metaInfo.DeleteFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}
