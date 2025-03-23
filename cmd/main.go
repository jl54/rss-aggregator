package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jl54/rss-aggregator/internal"
	"github.com/jl54/rss-aggregator/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable not defined")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL environment variable not defined")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to the database", err)
	}

	db := database.New(conn)
	apiCfg := internal.ApiConfig{
		DB: db,
	}

	go internal.StartScraping(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", internal.ReadyHandler)
	v1Router.Get("/err", internal.HandleError)
	v1Router.Post("/users", apiCfg.CreateUserHandler)
	v1Router.Get("/users", apiCfg.AuthMiddleware(apiCfg.GetUserByApiKeyHandler))
	v1Router.Post("/feeds", apiCfg.AuthMiddleware(apiCfg.CreateFeedHandler))
	v1Router.Get("/feeds", apiCfg.GetFeedsHandler)
	v1Router.Post("/feed-follows", apiCfg.AuthMiddleware(apiCfg.CreateFeedFollowHandler))
	v1Router.Get("/feed-follows", apiCfg.AuthMiddleware(apiCfg.GetFeedFollowsHandler))
	v1Router.Delete("/feed-follows/{feedFollowId}", apiCfg.AuthMiddleware(apiCfg.DeleteFeedFollowHandler))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %s", portString)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
