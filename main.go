package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// All global variables here
var (
	Db     *sql.DB
	RSAKey *rsa.PrivateKey
)

type App struct {
	Db     *sql.DB
	RSAKey *rsa.PrivateKey
}

// Initialize initializes database and rsa keypair
func (a *App) Initialize(config map[string]string) {
	// use Config map to connect to a database
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s",
		config["dbhost"], config["dbport"], config["dbname"], config["ssl"])
	var err error
	a.Db, err = sql.Open("postgres", connStr)

	// check if any errors during connection to database
	CheckErr(err)
	CheckErr(a.Db.Ping())

	// Get RSA keypair from file path
	a.RSAKey, err = GetKeyFromFile(config["rsakeypath"])
	CheckErr(err)
}

// WriteJSON returns a response in json format
func WriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(data)
	w.Write(response)
}

func main() {

	// initialize our application
	app := &App{}
	app.Initialize(Config)
	defer app.Db.Close()

	// handle routes to api
	router := mux.NewRouter()
	router.HandleFunc("/api/recipes", app.GetRecipes).Methods("GET")
	router.HandleFunc("/api/recipes/{id}", app.GetRecipeByID).Methods("GET")
	router.HandleFunc("/api/users/register", app.RegisterUser).Methods("POST")
	router.HandleFunc("/api/users/login", app.LoginUser).Methods("POST")
	router.HandleFunc("/api/users/auth", app.AuthUser).Methods("GET")

	// start server
	fmt.Println("Server listening on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// GetKeyFromFile retrieves private key from .pem file
func GetKeyFromFile(name string) (*rsa.PrivateKey, error) {
	// open .pem file
	file, err := os.Open(name)
	defer file.Close()
	CheckErr(err)

	// read content in .pem file
	var buffer bytes.Buffer
	buffer.ReadFrom(file)

	// decode the pem block and parse to a private key
	data, _ := pem.Decode(buffer.Bytes())
	return x509.ParsePKCS1PrivateKey(data.Bytes)
}

// CheckErr handles any errors that occurs
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
