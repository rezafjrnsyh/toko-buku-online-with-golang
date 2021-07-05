package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/url"
)

func ConnectDB() (*sql.DB, error) {
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "root", "", "127.0.0.1", "3306", "enigma_toko")
	val := url.Values{}
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	db, err := sql.Open(`mysql`, dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Print(err)
		_, _ = fmt.Scanln()
		log.Fatal(err)
	}
	log.Println("DataBase Successfully Connected")
	return db, err
}
