package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {

	dsName := "root:12345678@tcp(127.0.0.1:3306/po?charset=utf8)"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		fmt.Println(err)
	}
	// 最大空闲连接数
	db.SetMaxIdleConns(2)
	// 最大打开连接数
	db.SetMaxOpenConns(3)
	// 最大链接存活时间
	db.SetConnMaxIdleTime(7 * time.Hour)
	fmt.Println(db.Ping())
	defer db.Close()
}
