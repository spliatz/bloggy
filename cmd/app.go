package main

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/Intellect-Bloggy/bloggy-backend/internal/handler"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/repository"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/server"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/services"
)

func main() {

    dbPort, err := strconv.Atoi(os.Getenv("POSTGRES_LOCAL_PORT"))

    if dbPort > 1<<16 || dbPort < 1 || err != nil {
        logrus.Fatal("Подключение к базе данных не удалось: некорректный порт")
    }

    db, err := repository.NewPostgresDB(repository.PostgresConfig{
        Host:     os.Getenv("POSTGRES_HOST"),
        Port:     uint16(dbPort),
        Username: os.Getenv("POSTGRES_USER"),
        DBName:   os.Getenv("POSTGRES_DB"),
        Password: os.Getenv("POSTGRES_PASSWORD"),
    })

    if err != nil {
        logrus.Fatal("Подключение к базе данных не удалось: ", err)
    }

    logrus.Println("База данных успешно подключена")

    handlers := handler.NewHandlers(services.NewServices(repository.NewRepository(db)))
    srv := server.NewServer()

    srvPort, err := strconv.Atoi(os.Getenv("PORT"))
    if srvPort > 1<<16 || srvPort < 1 || err != nil {
        logrus.Fatal("Запуск сервера не удался: некорректный порт")
    }

    if err := srv.Run(uint16(srvPort), handlers.InitRoutes()); err != nil {
        logrus.Fatal("Ошибка запуска сервера: ", err)
    }
}
