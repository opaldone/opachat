// Package tools
package tools

import (
	"encoding/json"
	"os"
)

// Configuration config of the application
type Configuration struct {
	Appname  string              `json:"appname"`
	Address  string              `json:"address"`
	Port     int                 `json:"port"`
	Static   string              `json:"static"`
	Acme     bool                `json:"acme"`
	Acmehost []string            `json:"acmehost"`
	DirCache string              `json:"dirCache"`
	Crt      string              `json:"crt,omitempty"`
	Key      string              `json:"key,omitempty"`
	IceList  []map[string]string `json:"iceList"`
	Recorder *ConfRecorder       `json:"recorder,omitempty"`
}

type ConfRecorder struct {
	URLVirt  string `json:"urlVirt"`
	SoundLib string `json:"soundLib"`
	IHw      string `json:"iHw"`
	ScrRes   string `json:"scrRes"`
	LogLevel string `json:"logLevel"`
	Timeout  int    `json:"timeout"`
}

var (
	conf    *Configuration
	csrfkey string
)

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		Danger("Cannot open config.json file", err)
	}

	decoder := json.NewDecoder(file)
	conf = &Configuration{}
	err = decoder.Decode(conf)
	if err != nil {
		Danger("Cannot get configuration from file", err)
	}
}

func setCsrf() {
	csrfkey = CreateUUID()
}

// Env gets configuration
func Env(reload bool) *Configuration {
	if reload {
		loadConfig()
	}

	return conf
}

func GetIceList() ([]string, string, string) {
	urlsout := []string{}
	usernameout := ""
	credentialout := ""

	for _, v := range conf.IceList {
		urlsout = append(urlsout, v["urls"])

		if len(v["username"]) > 0 {
			usernameout = v["username"]
		}

		if len(v["credential"]) > 0 {
			credentialout = v["credential"]
		}
	}

	return urlsout, usernameout, credentialout
}

func GetKeyCSRF() string {
	return csrfkey
}
