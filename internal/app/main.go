package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"golang/internal/handler"
	"golang/internal/middleware"
	"golang/internal/repository"
	_postgres "golang/internal/repository/_postgres"
	"golang/internal/usecase"
	"golang/pkg/modules"

	"github.com/gorilla/mux"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreConfig()
	_postgre := _postgres.NewPGXDialect(ctx, dbConfig)
	repos := repository.NewRepositories(_postgre)
	userUsecase := usecase.NewUserUsecase(repos.UserRepository)
	userHandler := handler.NewUserHandler(userUsecase)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.AuthMiddleware)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "123123",
		DBName:      "mybd",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
