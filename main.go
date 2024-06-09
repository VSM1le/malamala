package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/VSM1le/malamala/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("port is not found in the enviroment")
	}
	portDB := os.Getenv("DB_URL")
	if portDB == "" {
		log.Fatal("DB url is not found in hte enviroment")
	}
	conn, err := sql.Open("postgres", portDB)
	if err != nil {
		log.Fatal("Can not connect to database:", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTION"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Post("/register", apiCfg.handlerRegisterUser)
	v1Router.Put("/login", apiCfg.handlerUserLogin)
	v1Router.Put("/logout", apiCfg.middlewareAuth(apiCfg.handlerUserLogout))

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server staring on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
