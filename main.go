package main

import (
	cors "github.com/rs/cors/wrapper/gin"
	"log"
	"main/config"
)

func main() {
	// fungsi untuk koneksi ke database
	db, err := config.ConnectDB()

	// apabila ada error program langsung berhenti
	if err != nil {
		log.Fatal(err.Error())
	}

	// Menginitialisasi sebuah object router
	r := config.CreateRouter()

	/*
		mekanisme untuk memberi tahu browser,
		apakah sebuah request yang di-dispatch dari aplikasi web domain lain atau origin lain,
		ke aplikasi web kita itu diperbolehkan atau tidak.
		Jika aplikasi kita tidak mengijinkan maka akan muncul error,
		dan request pasti digagalkan oleh browser. */
	r.Use(cors.AllowAll())

	/*
		Init Router digunakan untuk mengisi value Server struct
		Initial Routes digunakan untuk membuat router group
	*/
	config.InitRouter(db, r).InitializeRoutes()

	// Untuk menjalankan aplikasi
	errun := config.Run(r)
	if errun != nil {
		log.Fatal(errun.Error())
	}
}
