package db

import (
	"database/sql"
	"fmt"
	"log"
	mydb "myFileServer/db/mysql"
)

// OnFileUploadFinished 文件上传完成，保存meta
//
//	func OnFileUploadFinished(fileHash, fileName, fileAddr string, fileSize int64) bool {
//		stmt, err := mydb.DBConnect().Prepare(
//			"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`," +
//				"`file_addr`,`status`) values (?,?,?,?,1)")
//		if err != nil {
//			fmt.Println("Failed to prepare statement, err:" + err.Error())
//			return false
//		}
//		defer stmt.Close()
//
//		ret, err := stmt.Exec(fileHash, fileName, fileAddr, fileSize)
//		if err != nil {
//			fmt.Println(err.Error())
//			return false
//		}
//		if rf, err := ret.RowsAffected(); err == nil {
//			if rf <= 0 {
//				fmt.Printf("File with hash:%s has been uploaded before", fileHash)
//			}
//			return true
//		}
//		return false
//	}
func OnFileUploadFinished(fileHash, fileName, fileAddr string, fileSize int64) bool {
	db := mydb.DBConnect()
	if db == nil {
		fmt.Println("Failed to connect to database")
	}

	var version string
	err := db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatal("Failed to query database version: ", err)
	}
	fmt.Println("Connected to MySQL version: ", version)

	db, err = sql.Open("mysql", "reader:reader@tcp(192.168.2.238:3306)/fileserver?charset=utf8")
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}

	stmt, err := db.Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileHash, fileName, fileAddr, fileSize)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before", fileHash)
		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash           string
	FileName, FileAddr sql.NullString
	FileSize           sql.NullInt64
}

// GetFileMeta 从mysql获取文件元信息
func GetFileMeta(fileHash string) (*TableFile, error) {
	stmt, err := mydb.DBConnect().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file where file_sha1=? and status=1 limit 1",
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tFile := TableFile{}
	err = stmt.QueryRow(fileHash).Scan(&tFile.FileHash, &tFile.FileAddr, &tFile.FileName, &tFile.FileSize)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &tFile, nil
}
