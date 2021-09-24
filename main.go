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

	r := config.CreateRouter()

	config.InitRouter(db, r).InitializeRoutes()

	errun := config.Run(r)
	if errun != nil {
		log.Fatal(errun.Error())
	}
}
