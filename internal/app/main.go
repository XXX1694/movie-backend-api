package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "golang/docs"
	"golang/internal/cache"
	"golang/internal/handler"
	"golang/internal/middleware"
	"golang/internal/repository"
	_postgres "golang/internal/repository/_postgres"
	"golang/internal/usecase"
	"golang/internal/worker"
	"golang/pkg/modules"
)

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreConfig()
	_postgre := _postgres.NewPGXDialect(ctx, dbConfig)
	repos := repository.NewRepositories(_postgre)

	worker.StartUserCountWorker(repos)

	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisCache := cache.NewRedisCache(redisAddr)

	userUsecase := usecase.NewUserUsecase(repos.UserRepository, redisCache)
	userHandler := handler.NewUserHandler(userUsecase)

	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.AuthMiddleware)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	r.HandleFunc("/users/audit", userHandler.CreateUserWithAudit).Methods("POST")
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Server running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped gracefully")
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        getEnv("DB_HOST", "localhost"),
		Port:        getEnv("DB_PORT", "5432"),
		Username:    getEnv("DB_USERNAME", "postgres"),
		Password:    getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "mybd"),
		SSLMode:     getEnv("DB_SSLMODE", "disable"),
		ExecTimeout: 5 * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
