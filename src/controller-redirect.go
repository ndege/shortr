package main

import (
	_ "github.com/go-sql-driver/mysql"
  "github.com/gorilla/mux"
  "net/http"
)

// Handles a requested short URL.
// Redirects with a 301 header if found.
var RedirectController = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 1. Check if a slug exists
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// 2. Check if the slug exists in the database
	var url string
	err := db.QueryRow("SELECT `url` FROM `shortr` WHERE `slug` = ?", slug).Scan(&url)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 3. If the slug (and thus the URL) exist, update the hit counter
	stmt, err := db.Prepare("UPDATE `shortr` SET `hits` = `hits` + 1 WHERE `slug` = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Finally, redirect the user to the URL
	http.Redirect(w, r, url, http.StatusMovedPermanently)
})
