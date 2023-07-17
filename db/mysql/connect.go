package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "reader:reader@tcp(127.0.0.1:3306)/fileserver")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Printf("failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
}

//func init() {
//	var err error
//	db, err = sql.Open("mysql", "reader:reader@tcp(127.0.0.1:3306)/fileserver")
//	db.SetMaxOpenConns(1000)
//	err = db.Ping()
//	if err != nil {
//		fmt.Printf("failed to connect to mysql, err:" + err.Error())
//		os.Exit(1)
//	}
//}

// DBConnect 返回数据库链接对象
func DBConnect() *sql.DB {
	return db
}
