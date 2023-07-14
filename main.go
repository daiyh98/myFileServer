package main

import (
	"fmt"
	"myFileServer/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/succeed", handler.UploadSuccessHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server, error: %s", err.Error())
	}
}
