/* Shotr is an api to shorten url
*/
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Version of application.
const Version = "0.1.0"

var db  *sql.DB

func main() {

	cfgfile := flag.String("config", "env.json", "mandatory json config file")
	serve := flag.Bool("serve", false, "serve file")
	version := flag.Bool("version", false, "show version and exit")

	flag.Parse()

	if *version {
		log.Println(Version)
		os.Exit(0)
	}

	if *cfgfile == "" {
		log.Println("Config file to run service is mandatory. See env-example.json.")
		os.Exit(0)
	}
	ReadInConfig(*cfgfile)

  // Instantiate the database
	var err error
	dsn := cfg.DbUser + ":" + cfg.DbPass + "@tcp(" + cfg.DbHost + ":3306)/" + cfg.DbName + "?collation=utf8mb4_unicode_ci&parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if *serve {
		// Instantiate the mux router
		r := mux.NewRouter()
		r.Handle("/auth", AuthController).Methods("POST")
		r.Handle("/shortr", jwtMiddleware.Handler(GenerateController)).Methods("POST")
		r.Handle("/{slug:[a-z0-9]+}", RedirectController)
		r.Handle("/", IndexController)

		// Assign mux as the HTTP handler
		http.Handle("/", r)
		// Start HTTP Server
		log.Println("Start application v" + Version + " at port " + cfg.AppPort)
		err = http.ListenAndServe(":"+cfg.AppPort, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}

// Helper routine for sending JSON back to the client a bit more cleanly
func jResp(w http.ResponseWriter, data interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("Internal Server Error:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(payload)))
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})
