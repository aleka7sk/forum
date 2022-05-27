package server

import (
	"context"
	"database/sql"
	"forum/config"
	"forum/internal/auth"
	"forum/internal/post"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	authhttp "forum/internal/auth/delivery/http"
	authrepository "forum/internal/auth/repository/sqlite"
	authusecase "forum/internal/auth/usecase"
	posthttp "forum/internal/post/delivery/http"
	postrepository "forum/internal/post/repository"
	postusecase "forum/internal/post/usecase"
)

type App struct {
	httpServer  *http.Server
	Logger      *logrus.Logger
	AuthUseCase auth.UseCase
	PostUseCase post.UseCase
}

func NewApp(config config.Config) *App {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Db initialization error %v", err)
	}
	authrepository := authrepository.NewAuthRepository(db)
	postrepositry := postrepository.NewPostRepository(db)

	return &App{
		AuthUseCase: authusecase.NewService(authrepository, config.Hash_salt, []byte(config.Signing_key), config.Token_ttl),
		PostUseCase: postusecase.NewService(postrepositry),
		Logger:      logrus.New(),
	}
}

func (a *App) Run(config config.Config) error {
	router := http.NewServeMux()
	a.Logger.Info("Initialize router...")
	authhttp.RegisterHTTPEndpoints(router, a.AuthUseCase)
	posthttp.RegisterHTTPEndpoints(router, a.PostUseCase)
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
	// psqlInfo := fmt.Sprintf("host=%v port=%v user=%s "+
	// 	"password=%s sslmode=disable",
	// 	"localhost", "5432", "postgres", "password")
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	return nil, err
	// }

	// if err = db.Ping(); err != nil {
	// 	return nil, err
	// }

	// return db, nil
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	log.Print("Initialization of db")

	return db, nil
}
