package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strconv"
)

/*
creates config object
*/
func CreateConfig() *Config {
	return &Config{}
}

/*
creates user object
*/
func CreateUser() *User {
	return &User{}
}

/*
validates the config file, syntax, values and logic of config file
*/
func ValidateConfig(config *Config) bool {
	if !(config.Global.Authorization == "true" || config.Global.Authorization == "false") {
		config.logger("config incorrect : at global section directive authorization value must be true or false  ", 1)
		return false
	}
	if config.Global.Authorization == "true" {
		if len(config.Users) <= 0 {
			config.logger("config incorrect : in case of Authorization is true you have to define at least one user at Users section  ", 1)
			return false
		} else {
			for _, user := range config.Users {
				if user.Name == "" {
					config.logger("config incorrect : at users section user must have name ", 1)
					return false
				}
				if user.Token == "" {
					config.logger("config incorrect : at users section user must have token ", 1)
					return false
				}
				if user.Keys == "" {
					config.logger("config incorrect : at users section user must have keys ", 1)
					return false
				}
				if user.Cmds == "" {
					config.logger("config incorrect : at users section user must have cmds ", 1)
					return false
				}
			}
		}
	}
	if !(config.Cluster.Master == "true" || config.Cluster.Master == "false") {
		config.logger("config incorrect : at master section directive master value must be true or false  ", 1)
		return false
	}
	if config.Cluster.Master == "true" {
		if config.Cluster.CheckInterval <= 0 {
			config.logger("config incorrect : in case of master is true you have to define checkInterval (int) ", 1)
			return false
		}
		if len(config.Cluster.Slaves) <= 0 {
			config.logger("config incorrect : in case of master is true you have to define at least one slave at slaves section  ", 1)
			return false
		} else {
			for _, slave := range config.Cluster.Slaves {
				if !(slave.Id > 0) {
					config.logger("config incorrect : at slaves section you must set id >= 1 (int) for each slave ", 1)
					return false
				}
				if slave.Ip == "" {
					config.logger("config incorrect : at slaves section you must define ip for each slave ", 1)
					return false
				}
				if slave.Authorization == "true" {
					if slave.Token == "" {
						config.logger("config incorrect : at slaves section when authorization is true you must define token for each slave ", 1)
						return false
					}
				}
			}
		}
	}
	if !(config.Proc.Cores > 0) {
		config.logger("config incorrect : at proc section you must define cores ", 1)
		return false
	}
	if !(config.Proc.Threads > 0) {
		config.logger("config incorrect : at proc section you must define threads ", 1)
		return false
	}
	if config.Listen.Ip == "" {
		config.logger("config incorrect : at listen section you must define ip ", 1)
		return false
	}
	port, _ := strconv.Atoi(config.Listen.Port)
	if !(port > 1024) {
		config.logger("config incorrect : at listen section you must define port > 1024 ", 1)
		return false
	}
	if config.Listen.Protocol != "tcp" {
		config.logger("config incorrect : at listen section you must define protocol tcp ", 1)
		return false
	}
	if !(config.Log.Level < 3 && config.Log.Level >= 0) {
		config.logger("config incorrect : at log section level must be one of 0,1,2 ", 1)
		return false
	}

	return true
}

/*
reads the config file and return the config object
*/
func GetConfig(path string) *Config {
	file, err := os.Open(path)
	config := CreateConfig()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	if !ValidateConfig(config) {
		os.Exit(0)
	}
	config.logger("config.go - func GetConfig - with path : "+path, 2)
	config.logger("Nipo service is starting ...", 1)
	config.logger("Reading config file on :"+path, 1)
	return config
}

/*
reads the config file and return the config object
*/
func ReloadConfig(path string) (*Config, bool) {
	file, err := os.Open(path)
	config := CreateConfig()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	if !ValidateConfig(config) {
		return config, false
	}
	return config, true
}
