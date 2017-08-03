package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
)

/*
Config stores the basic configuration information for a CommComm server.
Port corresponds to the port the application will be running on. while DBInfo
defines all of the information needed to connect to the database used for storage of users and reports.
*/
type Config struct {
	Port   string `json:"port"`
	DB     DBInfo `json:"db"`
	Key    string `json:"key"`
	Cert   string `json:"cert"`
	Secret string `json:"secret"`
}

var conf Config

/*
DBInfo stores the information to connect to the database the CommComm server uses.
*/
type DBInfo struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Userpass string `json:"userpass"`
}

func main() {
	n := flag.String("config", "conf.json", "Configuration file. Must be JSON. Default is conf.json in the same working directory as the binary.")

	conf := getConfig(*n)
	err := InitDb(conf.DB.Username, conf.DB.Userpass, conf.DB.Address, conf.DB.Port)
	if err != nil {
		panic(err)
	}

	r := InitRouter()
	log.Fatal(http.ListenAndServe(":"+conf.Port, r))
}

func getConfig(name string) (c *Config) {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	d := json.NewDecoder(file)

	c = &Config{}
	err = d.Decode(&c)
	if err != nil {
		panic(err)
	}

	return
}
