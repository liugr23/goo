package database

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var SqlDB *sql.DB

func Init() {
	var err error
	SqlDB, err = sql.Open("mysql", "goo:Goo2017!@tcp(127.0.0.1:3306)/goo?parseTime=true")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = SqlDB.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}
