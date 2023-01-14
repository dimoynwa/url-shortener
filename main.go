package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/dimoynwa/url-shortener/api"
	"github.com/dimoynwa/url-shortener/repository/mongodb"
	"github.com/dimoynwa/url-shortener/repository/redis"
	"github.com/dimoynwa/url-shortener/shortener"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func chooseRepo() shortener.RedirectRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisUrl := os.Getenv("REDIS_URL")
		repo, err := redis.NewRedisRepository(redisUrl)
		if err != nil {
			log.Fatalf("Error creating Redis repository: %v\n", err)
		}
		return repo
	case "mongo":
		mongoUrl := os.Getenv("MONGO_URL")
		mongoDb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))

		repo, err := mongodb.NewMongoRepository(mongoUrl, mongoDb, mongoTimeout)
		if err != nil {
			log.Fatalf("Error creating Mongo repository: %v\n", err)
		}
		return repo
	default:
		panic("No Database configured")
	}
}

func main() {
	port := ":8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	fmt.Printf("Port : %v\n", port)

	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := api.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Printf("Listening on port %v\n", port)
		errs <- http.ListenAndServe(port, r)
	}()

	// If we clich CTRL + C, it will send signal to errs chan and will terminate the app
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf(":%s", <-c)
	}()

	fmt.Printf("Terminated: %v", <-errs)
}

// repo <- servive -> serializer -> http
