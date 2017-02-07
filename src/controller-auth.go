/* Get token
*/
package main

import (
  "crypto/sha256"
  "encoding/hex"
  "encoding/json"
  "github.com/dgrijalva/jwt-go"
  "net/http"
  "time"
)

var mySigningKey = []byte(cfg.SigningKey)

type AuthRequest struct {
    User			 string 	`json:"user"`
    Password	 string   `json:"password"`
}

type Token struct {
    Token    string    `json:"token"`
    Status 	 string    `json:"status"`
}

var AuthController = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

  // Get json POST request
	decoder := json.NewDecoder(r.Body)
	var param AuthRequest
  err := decoder.Decode(&param)
  if err != nil {
      jResp(w, Error{Msg: err.Error(), Status: "400"})
			return
  }

	// Check if the user parameter has been sent along (and is not empty)
	if param.User == "" {
		jResp(w, Error{Msg: "No parameter 'user' is set. Parameter is mandatory", Status: "400"})
		return
	}

  // Check if the password parameter has been sent along (and is not empty)
	if param.Password == "" {
		jResp(w, Error{Msg: "No parameter 'password' is set. Parameter is mandatory", Status: "400"})
		return
	}

  h := sha256.New()
  h.Write([]byte(param.Password))
  hashedPassword := hex.EncodeToString(h.Sum(nil))

  // Check if user already exists in the database with password
  var exists bool
  err = db.QueryRow("SELECT EXISTS(SELECT * FROM `users` WHERE `user` = ? AND `password` = ?)", param.User, hashedPassword).Scan(&exists)
  if err != nil {
		jResp(w, Error{Msg: err.Error(), Status: "400"})
		return
	}
  if exists == false {
    jResp(w, Error{Msg: "Unauthorized: Credentials are false.", Status: "401"})
    return
  }
  jResp(w, Token{Token: generateToken(), Status: "201"})
  return
})

func generateToken() string {

  /* Create the token */
  token := jwt.New(jwt.SigningMethodHS256)

  /* Create a map to store our claims */
  claims := token.Claims.(jwt.MapClaims)

  /* Set token claims */
  claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

  /* Sign the token with our secret */
  tokenString, _ := token.SignedString(mySigningKey)

  return string(tokenString);
}
