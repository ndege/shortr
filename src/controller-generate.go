package main

import (
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"errors"
  "net/http"
	"net/url"
  "math/rand"
	"strings"
	"log"
)

type Error struct {
    Msg      string    `json:"msg"`
    Status 	 string    `json:"status"`
}

type ShortrRequest struct {
		Url			 string 	`json:"url"`
	  Format	 string   `json:"wq"`
}

type Shorten struct {
    Url      string    `json:"url"`
    Status 	 string    `json:"status"`
}


// Shortens a given URL passed through in the request.
// If the URL has already been shortened, returns the existing URL.
// Writes the short URL in plain text to w.
<<<<<<< HEAD
func GenerateController(w http.ResponseWriter, r *http.Request) {
=======
var GenerateController = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
>>>>>>> 79a90de... add basic jwt authentification of api. add jwt create token at controller and jwt middleware to auth shortr post request

	// Get json POST request
	decoder := json.NewDecoder(r.Body)
	var param ShortrRequest
  err := decoder.Decode(&param)
  if err != nil {
      jResp(w, Error{Msg: err.Error(), Status: "400"})
			return
  }
	log.Println(param.Url)

	// Check if the url parameter has been sent along (and is not empty)
	if param.Url == "" {
		jResp(w, Error{Msg: "No parameter 'url' is set. Parameter is mandatory", Status: "400"})
		return
	}

	// Validate url
	err = validateUrlScheme(param.Url)
	if err  != nil {
		jResp(w, Error{Msg: err.Error(), Status: "400"})
		return
	}

	// Get the short URL out of the config
	if cfg.UrlService == "" {
		jResp(w, Error{Msg: "No base url: UrlService is set for service at configuration", Status: "400"})
		return
	}
	short_url := cfg.UrlService

	// Get ip address
	ip := httpRemoteIP(r)
	if throttleCheck(ip) == false {
		jResp(w, Error{Msg: "Limit achived of shorten ulr's in defined interval", Status: "400"})
		return
	}

	// Check if url already exists in the database
	var slug string
	err = db.QueryRow("SELECT `slug` FROM `shortr` WHERE `url` = ?", param.Url).Scan(&slug)
	if err == nil {
		// The URL already exists! Return the shortened URL.
    jResp(w, Shorten{Url: short_url + "/" + slug, Status: "201"})
		return
	}

	// generate a slug and validate it doesn't
	// exist until we find a valid one.
	var exists = true
	for exists == true {
		slug = generateSlug()
		err, exists = slugExists(slug)
		if err != nil {
			jResp(w, Error{Msg: err.Error(), Status: "400"})
			return
		}
	}

	// Insert it into the database
	stmt, err := db.Prepare("INSERT INTO `shortr` (`slug`, `url`, `date`, `hits`, `ip`) VALUES (?, ?, NOW(), ?, ?)")
	if err != nil {
		jResp(w, Error{Msg: err.Error(), Status: "400"})
		return
	}
	_, err = stmt.Exec(slug, param.Url, 0, ip)
	if err != nil {
		jResp(w, Error{Msg: err.Error(), Status: "400"})
		return
	}
	// Return response
  jResp(w, Shorten{Url: short_url + "/" + slug, Status: "201"})
	return

<<<<<<< HEAD
}
=======
})
>>>>>>> 79a90de... add basic jwt authentification of api. add jwt create token at controller and jwt middleware to auth shortr post request

// generateSlug will generate a random slug to be used as shorten link.
func generateSlug() string {
	// It doesn't exist! Generate a new slug for it
	// From: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
	var chars = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	s := make([]rune, 6)
	for i := range s {
		s[i] = chars[rand.Intn(len(chars))]
	}

	return string(s)
}

func httpRemoteIP(r *http.Request) string {
	idx := strings.LastIndex(r.RemoteAddr, ":")
	if idx == -1 {
		return r.RemoteAddr
	}
	return r.RemoteAddr[:idx]
}

// slugExists will check whether the slug already exists in the database
func slugExists(slug string) (e error, exists bool) {
	err := db.QueryRow("SELECT EXISTS(SELECT * FROM `shortr` WHERE `slug` = ?)", slug).Scan(&exists)
	if err != nil {
		return err, false
	}

	return nil, exists
}

func throttleCheck(ip string) bool {
	rows, err := db.Query("SELECT COUNT(*) as count FROM `shortr` WHERE `ip` = ? AND date > CURRENT_TIMESTAMP - INTERVAL 1 hour", ip)
	if err != nil {
		return false
	}
	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return false
		}
	}
	return count < cfg.MaxRequest
}

func validateUrlScheme(urlf string) error {
	urlValidate, err := url.Parse(urlf)
	if err != nil {
		return errors.New("Cannot parse url")
	}
	// We expect a URL has at least one period.
	if !strings.Contains(urlValidate.Host, ".") {
		return errors.New("Invalid url")
	}
	baseUrl, err := url.Parse(cfg.UrlService)
	if err != nil {
		return errors.New("Cannot parse baseUrl")
	}
	// To prevent someone from building predictive redirect chains to try and overload us
	if urlValidate.Host == baseUrl.Host {
		return errors.New("Short urls pointing to "+cfg.UrlService+" are not allowed")
	}
	return nil
}
