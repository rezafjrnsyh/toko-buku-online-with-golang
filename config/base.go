package config

import (
	"database/sql"
	"fmt"
	"log"
	"main/controllers"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

const MODE = "release"

type Server struct {
	DB	*sql.DB
	Router *gin.Engine
}

func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func ConnectDB() (*sql.DB, error) {
	// get := GetEnvWithKey
	// DB_USER := get("DB_USER")
	// DB_PASS := get("DB_PASS")
	// DB_HOST := get("DB_HOST")
	// DB_PORT := get("DB_PORT")
	// DB_NAME := get("DB_NAME")
	// DB_LOC := get("DB_LOC")
	DB_USER := viper.GetString("database.DB_USER")
	DB_PASS := viper.GetString("database.DB_PASS")
	DB_HOST := viper.GetString("database.DB_HOST")
	DB_PORT := viper.GetString("database.DB_PORT")
	DB_NAME := viper.GetString("database.DB_NAME")
	DB_LOC := viper.GetString("database.DB_LOC")
	//String format untuk koneksi
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)
	val := url.Values{}
	// menambahkan value location
	val.Add("loc", DB_LOC)
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	// Buka koneksi database
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

func CreateRouter() *gin.Engine {

	// Untuk Set Mode
	if  MODE != "release" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Menginstance object router
	r := gin.Default()
	return r
}

func InitRouter(db *sql.DB, r *gin.Engine) *Server {
	return &Server{
		DB: db,
		Router: r,
	}
}

func (server *Server) InitializeRoutes()  {
	// Membuat sebuah group router
	r := server.Router.Group("v1")
	controllers.NewCategoryController(server.DB, r)
	controllers.NewBookController(server.DB, r)
	controllers.NewMemberController(server.DB, r)
}

func Run(r *gin.Engine) error {
	fmt.Println("Listening to port 8801")
	err := r.Run(":8801")
	return err
}
