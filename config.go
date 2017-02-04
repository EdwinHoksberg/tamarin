package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Hostname      string
	Listen        string
	Port          int
	Read_timeout  int
	Write_timeout int
	Loglevel      string
}

func (c *Config) loadConfig(configFile string) bool {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Could not load config file: %s", err.Error())
		return false
	}

	err_ := json.Unmarshal(file, &c)
	if err_ != nil {
		fmt.Printf("Could not parse json: %s", err_.Error())
		return false
	}

	return true
}
