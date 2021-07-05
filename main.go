package main

import (
	"log"
	"main/config"
)

func main()  {
	// fungsi untuk koneksi ke database
	db,err := config.ConnectDB()

	// apabila ada error program langsung berhenti
	if err !=nil {
		log.Fatal(err.Error())
	}
}
