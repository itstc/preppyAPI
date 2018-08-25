package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"golang.org/x/crypto/bcrypt"
)

// JSONAuthError is the standard response for errors in authenticating
var (
	JSONAuthError = map[string]interface{}{
		"error": "invalid authentication!",
		"auth":  false,
	}
	JSONLoginError = map[string]interface{}{
		"error": "invalid login!",
	}
	JSONRegisterError = map[string]interface{}{
		"error": "unable to process registration!",
	}
)

// AuthUser is a get method that returns if token is valid
func (a *App) AuthUser(w http.ResponseWriter, r *http.Request) {
	// logging request
	fmt.Printf("%s > GET Auth User\n", r.RemoteAddr)

	jwt, err := jws.ParseJWTFromRequest(r)
	// no token found in request
	if err != nil {
		WriteJSON(w, JSONAuthError)
		return
	}

	// check if token is valid
	if err := jwt.Validate(a.RSAKey.Public(), crypto.SigningMethodRS256); err != nil {
		WriteJSON(w, JSONAuthError)
		return
	}

	// successfully verified token
	WriteJSON(w, map[string]interface{}{
		"message": "auth successful!",
		"auth":    true,
	})
}

// RegisterUser takes email, name, and password from body and creates a new user in database
func (a *App) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// logging request
	fmt.Printf("%s > POST Register User\n", r.RemoteAddr)

	var jr map[string]interface{}

	// Read and proccess body request
	var buffer bytes.Buffer
	buffer.ReadFrom(r.Body)
	err := json.Unmarshal(buffer.Bytes(), &jr)

	// error occurred proccessing form data
	if err != nil {
		WriteJSON(w, JSONRegisterError)
		return
	}

	// 1 uppercase and password length >= 8
	// if match, _ := regexp.MatchString("^(?=.*[A-Z]){8,}$", jr["password"].(string)); !match {
	// 	responseEncoder.Encode(map[string]string{"error": "Invalid Password! (1 uppercase and at least 8 characters)"})
	// 	return
	// }

	// convert string password to hash
	jr["password"], err = bcrypt.GenerateFromPassword([]byte(jr["password"].(string)), HASHCOST)
	if err != nil {
		WriteJSON(w, JSONRegisterError)
		return
	}

	// try inserting into database
	_, err = a.Db.Exec("INSERT INTO account(email,name,password) VALUES ($1, $2, $3)", jr["email"], jr["name"], jr["password"])
	if err != nil {
		WriteJSON(w, JSONRegisterError)
		return
	}

	// everything went well return success response
	WriteJSON(w, map[string]string{"ok": "1", "message": "user successfully created!"})
}

// LoginUser takes email and password from body and returns a authentication token if successful
func (a *App) LoginUser(w http.ResponseWriter, r *http.Request) {
	// logging request
	fmt.Printf("%s > POST Login User\n", r.RemoteAddr)

	var jr map[string]interface{}

	// Read and proccess body request
	var buffer bytes.Buffer
	buffer.ReadFrom(r.Body)
	err := json.Unmarshal(buffer.Bytes(), &jr)
	if err != nil {
		WriteJSON(w, JSONLoginError)
		return
	}

	// query for the hash password given email
	var hashPassword []byte
	var userID int
	var username string
	row := a.Db.QueryRow("SELECT id, name, password FROM account WHERE email = $1", jr["email"].(string))
	err = row.Scan(&userID, &username, &hashPassword)

	// error querying for password
	if err != nil {
		WriteJSON(w, JSONLoginError)
		return
	}

	// no password found in query (account DNE)
	if len(hashPassword) == 0 {
		WriteJSON(w, JSONLoginError)
		return
	}

	// error when comparing hash and password
	err = bcrypt.CompareHashAndPassword(hashPassword, []byte(jr["password"].(string)))
	if err != nil {
		WriteJSON(w, JSONLoginError)
		return
	}

	var claims = jws.Claims{
		"id":    userID,
		"name":  username,
		"email": jr["email"].(string),
	}

	jwt := jws.NewJWT(claims, crypto.SigningMethodRS256)
	token, err := jwt.Serialize(a.RSAKey)

	if err != nil {
		WriteJSON(w, JSONLoginError)
		return
	}

	// password matches hashed password so login successful
	WriteJSON(w, map[string]string{
		"ok":      "1",
		"message": "successfully logged in!",
		"token":   string(token),
	})
}
