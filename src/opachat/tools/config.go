package tools

import (
	"encoding/json"
	"os"
)

// Configuration config of the application
type Configuration struct {
	Appname      string              `json:"appname"`
	Address      string              `json:"address"`
	Debug        bool                `json:"debug"`
	Acmehost     string              `json:"acmehost"`
	DirCache     string              `json:"dirCache"`
	Port         int                 `json:"port"`
	Crt          string              `json:"crt,omitempty"`
	Key          string              `json:"key,omitempty"`
	ReadTimeout  int64               `json:"readTimeout"`
	WriteTimeout int64               `json:"writeTimeout"`
	IdleTimeout  int64               `json:"idleTimeout"`
	IceList      []map[string]string `json:"iceList"`
	Saver        *ConfSaver          `json:"saver"`
}

type ConfSaver struct {
	UVirt    string `json:"u_virt"`
	Loop     string `json:"loop"`
	Screen   string `json:"screen"`
	Loglevel string `json:"loglevel"`
}

var conf *Configuration
var csrf_key string

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

	csrf_key = CreateUUID()
}

// Env gets configuration
func Env() *Configuration {
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
