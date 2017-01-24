package main

import (
	_ "github.com/go-sql-driver/mysql"
	"errors"
  "net/http"
	"net/url"
  "math/rand"
	"strings"
)

type Error struct {
    Msg      string    `json:"msg"`
    Status 	 string    `json:"status"`
}

type Shorten struct {
    Url      string    `json:"url"`
    Status 	 string    `json:"status"`
}


// Shortens a given URL passed through in the request.
// If the URL has already been shortened, returns the existing URL.
// Writes the short URL in plain text to w.
func GenerateController(w http.ResponseWriter, r *http.Request) {

  // Get back as json format
  output := r.URL.Query().Get("output")
  json := false
  if output == "json" {
    json = true;
  }

	// Check if the url parameter has been sent along (and is not empty)
	url := r.URL.Query().Get("url")
  if url == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Validate url
	err := validateUrlScheme(url)
	if err != nil {
		if (json) {
			jResp(w, Error{Msg: err.Error(), Status: "400"})
		} else {
			http.Error(w, "", http.StatusBadRequest)
		}
		return
	}

	// Get the short URL out of the config
	if cfg.UrlService == "" {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	short_url := cfg.UrlService

	// Get ip address
	ip := httpRemoteIP(r)
	if throttleCheck(ip) == false {
		if (json) {
				jResp(w, Error{Msg: "Limit achived of shorten ulr's in defined interval", Status: "400"} )
		} else {
			  http.Error(w, "", http.StatusBadRequest)
		}
		return
	}

	// Check if url already exists in the database
	var slug string
	err = db.QueryRow("SELECT `slug` FROM `shortr` WHERE `url` = ?", url).Scan(&slug)
	if err == nil {
		// The URL already exists! Return the shortened URL.
    if json {
        jResp(w, Shorten{Url: short_url + "/" + slug, Status: "200"} )
    } else {
		    w.Write([]byte(short_url + "/" + slug))
    }
    return
	}

	// generate a slug and validate it doesn't
	// exist until we find a valid one.
	var exists = true
	for exists == true {
		slug = generateSlug()
		err, exists = slugExists(slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Insert it into the database
	stmt, err := db.Prepare("INSERT INTO `shortr` (`slug`, `url`, `date`, `hits`, `ip`) VALUES (?, ?, NOW(), ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(slug, url, 0, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

  if json {
      jResp(w, Shorten{Url: short_url + "/" + slug, Status: "200"})
  } else {
	   w.WriteHeader(http.StatusCreated)
	   w.Write([]byte(short_url + "/" + slug))
  }
}

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
