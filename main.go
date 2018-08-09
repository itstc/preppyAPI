package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Db is our connection to our database
var Db *sql.DB

func main() {
	// initialize database
	initDB()
	defer Db.Close()

	// handle routes to api
	router := mux.NewRouter()
	router.HandleFunc("/api/recipes", GetRecipes).Methods("GET")
	router.HandleFunc("/api/recipes/{id}", GetRecipeByID).Methods("GET")

	// start server
	fmt.Println("Server listening on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func initDB() {
	// use Config map to connect to a database
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s",
		Config["dbhost"], Config["dbport"], Config["dbname"], Config["ssl"])
	var err error
	Db, err = sql.Open("postgres", connStr)

	// check if any errors during connection to database
	checkErr(err)
	checkErr(Db.Ping())
}

// checkErr handles any errors that occurs
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
