package config

import (
	"fmt"
	"log"
	"main/controllers"
	"os"

	"github.com/gin-gonic/gin"
	// _ "github.com/jackc/pgx/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const MODE = "release"

type Server struct {
	DB	*sqlx.DB
	Router *gin.Engine
}

type DbConn struct {
	Db *gorm.DB
}

func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func NewDbConn() *DbConn {
	dbHost := "localhost"
	dbPort := "5432"
	dbUser := "reza"
	dbPassword := "1234"
	dbName := "toko_buku"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)

	env := "dev"
	//Daftar konfigurasi GORM apa saja yang ada
	//https://gorm.io/docs/gorm_config.html
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	if env == "dev" {
		return &DbConn{
			Db: db.Debug(),
		}
	}

	return &DbConn{
		Db: db,
	}

}

func ConnectDB() (*sqlx.DB, error) {
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
	// DB_LOC := viper.GetString("database.DB_LOC")
	//String format untuk koneksi
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", DB_HOST, DB_USER, DB_NAME, DB_PASS, DB_PORT)
	// val := url.Values{}
	// menambahkan value location
	// val.Add("loc", DB_LOC)
	// dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	// Buka koneksi database
	db, err := sqlx.Connect(`pgx`, dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Print(err)
		_, _ = fmt.Scanln()
		log.Fatal(err)
	}
	log.Println("Database Successfully Connected")
	return db, err
}

func (d *DbConn) Migration(tables ...interface{}) {
	err := d.Db.AutoMigrate(tables...)
	if err != nil {
		log.Fatalln(err)
	}
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

func InitRouter(db *sqlx.DB, r *gin.Engine) *Server {
	return &Server{
		DB: db,
		Router: r,
	}
}

func (server *Server) InitializeRoutes()  {
	// Membuat sebuah group router
	r := server.Router
	// controllers.NewCategoryController(server.DB, r)
	controllers.NewBookController(server.DB, r)
	controllers.NewMemberController(server.DB, r)
}

func Run(r *gin.Engine) error {
	fmt.Println("Listening to port 8801")
	err := r.Run(":8801")
	return err
}
