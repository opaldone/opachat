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
	UrlVirt  string `json:"urlVirt"`
	SoundLib string `json:"soundLib"`
	IHw      string `json:"iHw"`
	ScrRes   string `json:"scrRes"`
	LogLevel string `json:"logLevel"`
	Timeout  int    `json:"timeout"`
}

var (
	conf     *Configuration
	csrf_key string
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
	csrf_key = CreateUUID()
}

// Env gets configuration
func Env(reload bool) *Configuration {
	if reload {
		loadConfig()
	}

	return conf
}

func GetIceList() ([]string, string, string) {
	urls_out := []string{}
	username_out := ""
	credential_out := ""

	for _, v := range conf.IceList {
		urls_out = append(urls_out, v["urls"])

		if len(v["username"]) > 0 {
			username_out = v["username"]
		}

		if len(v["credential"]) > 0 {
			credential_out = v["credential"]
		}
	}

	return urls_out, username_out, credential_out
}

func GetKeyCSRF() string {
	return csrf_key
}
