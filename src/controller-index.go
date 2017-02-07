package main

import (
  "net/http"
)

// Catches all other requests to the short URL domain.
// If a default URL exists in the config redirect to it.
var IndexController = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// 1. Get the redirect URL out of the config
	if cfg.UrlFallback == "" {
		// The reason for using StatusNotFound here instead of StatusInternalServerError
		// is because this is a catch-all function. You could come here via various
		// ways, so showing a StatusNotFound is friendlier than saying there's an
		// error (i.e. the configuration is missing)
		http.NotFound(w, r)
		return
	}

	// 2. If it exists, redirect the user to it
	http.Redirect(w, r, cfg.UrlFallback, http.StatusMovedPermanently)
})
