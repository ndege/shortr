// All hard-coded app configuration is in this file, as is all code for
// interacting with config information that is stored in a JSON file.
package main

import (
	"encoding/json"
	"log"
	"os"
)

// Configuration information MySQL.  We parse this from a JSON config file

// NB: field names must start with Capital letter for JSON parse to work
// NB: between the field names and JSON mnemonics, it should be easy to
//     figure out what each field does
type Config struct {
	DbUser       string   `json:"DbUsername"`
	DbPass       string   `json:"DbPassword"`
	DbHost       string   `json:"DbHost"`
	DbPort       string   `json:"DbPort"`
	DbName       string   `json:"DbName"`
  UrlService   string   `json:"ShortUrl`
	UrlFallback  string   `json:"FallbackUrl`
  AppPort      string   `json:"ApplicationPort"`
	MaxRequest   int      `json:"MaximalRequest"`
	SigningKey   string   `json:"JwtSigningKey"`
}

// The configuration information for the app we're administering
var cfg Config

// Load a JSON file that has all the config information for our app, and put
// the JSON contents into the cfg variable
func ReadInConfig(cfgFileName string) {
	// first, load the JSON file and parse it into /cfg/
	f, err := os.Open(cfgFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	jsonParser := json.NewDecoder(f)
	if err = jsonParser.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
}
