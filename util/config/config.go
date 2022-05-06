package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	Config = make(map[string]interface{})
)

func LoadConfig(path string) error {
	confFile, err := os.Open(path)
	if err != nil {
		log.Println("Open config file err:", err)
		return err
	}
	by, err := ioutil.ReadAll(confFile)
	if err != nil {
		log.Println("Read config err:", err)
		return err
	}
	err = json.Unmarshal(by, &Config)
	if err != nil {
		log.Println("json.Unmarshal config error:", err)
		return err
	}
	return nil
}