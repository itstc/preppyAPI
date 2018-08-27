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

	"github.com/go-redis/redis"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Db     *sql.DB
	RSAKey *rsa.PrivateKey
	Redis  *redis.Client
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

	a.Redis = redis.NewClient(&redis.Options{
		Addr:     Config["redishost"],
		Password: Config["redispassword"],
		DB:       0, // use default DB
	})

	_, err = a.Redis.Ping().Result()
	CheckErr(err)

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

	router.Handle("/api/users/register", app.Auth(http.HandlerFunc(app.RegisterUser))).Methods("POST")
	router.Handle("/api/users/login", app.Auth(http.HandlerFunc(app.LoginUser))).Methods("POST")
	router.Handle("/api/users/logout", app.Auth(http.HandlerFunc(app.LogoutUser))).Methods("POST")

	router.HandleFunc("/api/users/auth", app.AuthUser).Methods("GET")

	// global middlewares
	router.Use(LoggingMiddleware)

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
