package app

import (
	"os"
	"valerii/crudbananas/internal/delivery/rest"
	"valerii/crudbananas/internal/repository/pdb"
	"valerii/crudbananas/internal/server"
	"valerii/crudbananas/internal/service"
	"valerii/crudbananas/pkg/database"

	"os/signal"
	"syscall"

	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

func Run() {
	initLogger()

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env file: %s", err.Error())
	}

	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := pdb.NewRepository(db)
	services := service.NewService(repos)
	handlers := rest.NewHandler(services)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()
	log.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Errorf("Error closing database connection: %s", err.Error())
	}

	log.Info("Server exited gracefully")
}

func initLogger() {
	log.SetFormatter(new(log.JSONFormatter))
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
