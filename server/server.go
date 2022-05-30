package server

import (
	"context"
	"database/sql"
	"fmt"
	"forum/config"
	"forum/internal/auth"
	"forum/internal/post"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	authhttp "forum/internal/auth/delivery/http"
	authrepository "forum/internal/auth/repository/sqlite"
	authusecase "forum/internal/auth/usecase"
	posthttp "forum/internal/post/delivery/http"
	postrepository "forum/internal/post/repository"
	postusecase "forum/internal/post/usecase"

	"github.com/go-redis/redis"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type App struct {
	httpServer  *http.Server
	Logger      *logrus.Logger
	AuthUseCase auth.UseCase
	PostUseCase post.UseCase
	Db          *sql.DB
	Redis       *redis.Client
}

func NewApp(config config.Config) *App {
	db, err := initDB()
	if err != nil {
		log.Fatalf("SQL Db initialization error %v", err)
	}
	redis, err := initRedis()
	if err != nil {
		log.Fatalf("Redis Db initialization error %v", err)
	}
	authrepository := authrepository.NewAuthRepository(db, redis)
	postrepository := postrepository.NewPostRepository(db)

	return &App{
		AuthUseCase: authusecase.NewService(authrepository, config.Hash_salt, []byte(config.Signing_key), config.Token_ttl),
		PostUseCase: postusecase.NewService(postrepository, config.Hash_salt, []byte(config.Signing_key), config.Token_ttl),
		Logger:      logrus.New(),
		Db:          db,
		Redis:       redis,
	}
}

func (a *App) Run(config config.Config) error {
	router := http.NewServeMux()
	defer a.Db.Close()
	defer a.Redis.Close()
	a.Logger.Info("Initialize router...")
	authhttp.RegisterHTTPEndpoints(router, a.AuthUseCase, a.Redis)
	posthttp.RegisterHTTPEndpoints(router, a.PostUseCase, a.Redis)

	a.Logger.Info("Register HTTP endpoints...")

	a.httpServer = &http.Server{
		Addr:           ":" + config.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			a.Logger.Fatalf("Failed to listen and server: %+v", err)
		}
	}()
	a.Logger.Printf("Server run on port: %v", config.Port)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()
	return a.httpServer.Shutdown(ctx)
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	log.Print("Initialization of db")
	return db, nil
}

func initRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Printf("Redis initialization error: %v", err)
		return nil, err
	}
	fmt.Println(pong)
	return client, nil
}
