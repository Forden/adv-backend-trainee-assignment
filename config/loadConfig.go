package config

import (
	"encoding/json"
	"os"
)

type MyConfig struct {
	Port       int    `json:"port"`
	Host       string `json:"host"`
	UsedDB     string `json:"used_db"`
	PostgreSQL struct {
		Port        int    `json:"port"`
		Host        string `json:"host"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		DBName      string `json:"db_name"`
		SSLMode     string `json:"ssl_mode"`
		SSLRootCert string `json:"ssl_root_cert"`
	} `json:"postgresql"`
}

func LoadConfig(filename string) (MyConfig, error) {
	var cfg MyConfig
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&cfg)
	if err != nil {
		return MyConfig{}, err
	} else {
		return cfg, nil
	}
}
