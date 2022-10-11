package main

import (
	"github.com/Intellect-Bloggy/bloggy-backend/internal/handler"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/repository"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/server"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/services"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func main() {

	port, err := strconv.Atoi(os.Getenv("POSTGRES_LOCAL_PORT"))

	if port > 1<<16 || port < 1 || err != nil {
		logrus.Fatal("Подключение к базе данных не удалось: некорректный порт")
	}

	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     uint16(port),
		Username: os.Getenv("POSTGRES_USER"),
		DBName:   os.Getenv("POSTGRES_DB"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	})

	if err != nil {
		logrus.Fatal("Подключение к базе данных не удалось")
	}

	logrus.Println("База данных успешно подключена")

	handlers := handler.NewHandlers(services.NewServices(repository.NewRepository(db)))
	srv := server.NewServer()

	if err := srv.Run(os.Getenv("PORT"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Server start error: %s", err.Error())
	}
}
