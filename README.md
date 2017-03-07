# An URL shortener written in Golang
Shortr is inspired by Sam Wierema' [Go URL Shortener](https://github.com/samwierema/go-url-shortener), and more a project to get familiar with Golang.

## Features

Scope of the application is to shorten urls using only `[a-z0-9]` characters and
redirect requests to the corresponding site. In addtion the application has a
tiny api interface to auhenticate and create short urls with a jwt token. Use
case is the idea of a user based url shortener api due to common laws of countries
to be responsible for content of provided services.

In addition, ther are several features implemented as:
* Redirect to your main website when no slug, or incorrect slug, is entered, e.g. `http://domain.tdl/` → `http://website.domain.tdl/`.
* Doesn’t create a short URLs again if there's an attempt to shorten same URL. Therefor script returns already existing short URL.
* Additionally validation and security checks as: (1) Avoid flooding. Limit creation of short urls in a defined time interval. (2) Check if url host is valid. (3) Avoiding self reference on base url.

### API Requests

All requests json-encoded and returns as response a json.

| Requests   | Variables                                 | Type   | Response  															| Token
|------------|-------------------------------------------|--------|-----------------------------------------| ------
| /auth      | {'username':{user},'password':{password}} | POST   | {'url':{shortr_url},'status':{2xx}}     |
| /shortr    | {'url':{url_to_shorten}}                  | POST   | {'token':{bearer_token},'status':{2xx}} | X

Please note error response will return {'error':{error_msg},'status':{4xx}}

_Examples:_
curl -X POST "domain.tdl/auth" -H "Content-Type: application/json" -d "{\"user\":\"test\",\"password\":\"pass\"}"
curl -X POST "domain.tdl/shortr" -H "Content-Type: application/json" -H "Authorization: bearer {token}" -d "{\"url\":\"domain_to_shorten.tdl\"}"

## Installation
1. Download the source code and install it using the `go install` command.
2. Use `database.sql` in `install/db` to create tables.
3. Create a config file in `/path/to/shortr/` named `env.json`. Use `env-example.json` as a example.
4. Run the program as a daemon using one of the many methods: write a script for [upstart](https://launchpad.net/upstart), init, use [daemonize](http://software.clapper.org/daemonize/), [Supervisord](http://supervisord.org/), [Circus](http://circus.readthedocs.org/) or just plain old `nohup`. You can even start (and manage) it in a `screen` session.
5. Adding the following configuration to Apache (make sure you've got [mod_proxy](http://httpd.apache.org/docs/2.2/mod/mod_proxy.html) enabled):
```
<VirtualHost *:80>
	ServerName your-short-domain.ext

	ProxyPreserveHost on
	ProxyPass / http://localhost:8080/
	ProxyPassReverse / http://localhost:8080/
</VirtualHost>
```

### Using the example init script
You will find an example init script in the `scripts` folder. To use, you **must** at least change the GOPATH line to point to your Go root path.
