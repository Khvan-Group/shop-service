package main

import (
	"github.com/Khvan-Group/common-library/logger"
	"github.com/Khvan-Group/common-library/utils"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"net/http"
	"shop-service/internal/api"
	"shop-service/internal/core/rabbitmq"
	"shop-service/internal/db"
)

const SERVER_PORT = "SERVER_PORT"

func main() {
	start()
}

func start() {
	// init logger
	logger.InitLogger()
	logger.Logger.Info("Starting server")

	// load environments
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	port := ":" + utils.GetEnv(SERVER_PORT)

	// init RabbitMQ
	rabbitmq.InitRabbitMQ()

	// init DB
	db.InitDB()
	r := mux.NewRouter()
	srv := api.New()
	srv.AddRoutes(r)

	logger.Logger.Fatal(http.ListenAndServe(port, r).Error())
}
