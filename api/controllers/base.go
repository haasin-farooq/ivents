package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres

	"github.com/haasin-farooq/ivents/api/middlewares"
	"github.com/haasin-farooq/ivents/api/models"
	"github.com/haasin-farooq/ivents/api/responses"
)
type App struct {
	Router *mux.Router
	DB *gorm.DB
}

// Initialize connect to the database and wire up routes
func (a *App) Initialize(DbHost, DbPort, DbUser, DbName, DbPassword string) {
	var err error

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	a.DB, err = gorm.Open("postgres", connectionString)
	if err != nil {
		fmt.Printf("\nCannot connect to database %s\n", DbName)
		log.Fatal("This is the error:", err)
	}

	a.DB.Debug().AutoMigrate(&models.User{})

	a.Router = mux.NewRouter().StrictSlash(true)
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.Use(middlewares.SetContentTypeMiddleware) // setting content-type to json

	a.Router.HandleFunc("/", home).Methods("GET")
	a.Router.HandleFunc("/register", a.UserSignUp).Methods("POST")
	a.Router.HandleFunc("/login", a.Login).Methods("POST")
}

func (a *App) RunServer() {
	log.Printf("\nServer starting on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To Ivents")
}