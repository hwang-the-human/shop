package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	orderscontrollers "shop/apps/orders/infrastructure/controllers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	ordersservices "shop/apps/orders/application/usecases"
	ordersrepository "shop/apps/orders/infrastructure/repository"
	"shop/shared/migrations"
	sharedrepository "shop/shared/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	db, err := sharedrepository.NewPostgresRepository(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}()

	if err := migrations.ApplyMigrations(db.DB, "./apps/orders/infrastructure/migrations"); err != nil {
		log.Fatal(err)
	}

	ordersRepo := ordersrepository.NewOrdersRepository(db)
	ordersService := ordersservices.NewOrdersService(ordersRepo)
	ordersController := orderscontrollers.NewOrderController(ordersService)

	r := mux.NewRouter()

	r.HandleFunc("/orders", ordersController.CreateOrder).Methods("POST")
	r.HandleFunc("/orders/{id}", ordersController.GetOrder).Methods("GET")
	r.HandleFunc("/orders/{id}", ordersController.UpdateOrder).Methods("PUT")
	r.HandleFunc("/orders/{id}", ordersController.DeleteOrder).Methods("DELETE")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

//Спросить про три слоя репо заказов, репо бд и сервисы, пагинацию
