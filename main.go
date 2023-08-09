package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/abhishekghimire40/blog-feed-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbQueries := database.New(db)

	dbConfig := apiConfig{
		DB: dbQueries,
	}

	// main router
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// version router: routers staring with /v1
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleReadiness())
	v1Router.Get("/err", handleError())
	// users endpoints
	v1Router.Get("/users", dbConfig.middlewareAuth(dbConfig.GetUser))
	v1Router.Post("/users", createUser(dbConfig.DB))
	// feeds endpoints
	v1Router.Get("/feeds", getAllFeeds(dbConfig.DB))
	v1Router.Post("/feeds", dbConfig.middlewareAuth(dbConfig.createFeeds))
	router.Mount("/v1", v1Router)

	// starting the server and serving at port
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	fmt.Println("serving on port:", port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server")
	}
}
