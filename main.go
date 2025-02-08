package main

import (
	"io"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
)

type User struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
  Username  string `json:"username"`
	Password string `json:"password"`
}

var DB *gorm.DB

func ConnectDatabase() {

  log.Println("Connecting to db...")

	dsn := "host=HOST user=USER password=PASSWORD dbname=DB port=5432 sslmode=SSLMODE"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	DB = db
}

func main() {

  ConnectDatabase()

  log.Println("Starting server...")
  
  r := chi.NewRouter()
  r.Use(middleware.Logger)

  log.Println("Server started on PORT")

  r.Get("/", GetUsers)

	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./openapi.yaml") 
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/openapi.yaml"), 
	))

  r.Post("/send", CreateUser)

  srv := http.Server{
    Addr: PORT,
    Handler: r,
  }
  srv.ListenAndServe()
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	var users []User
	DB.Find(&users)

  bytes, _ := sonic.Marshal(users)

  w.Write(bytes)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

  body, _ := io.ReadAll(r.Body)

  var user User
  sonic.Unmarshal(body, &user)

  user.ID = uuid.New()
	DB.Create(&user)

  bytes, _ := sonic.Marshal(user)

  w.Write(bytes)
}
