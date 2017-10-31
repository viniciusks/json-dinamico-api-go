package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	MongoDB struct {
		Host   string `json:"host"`
		DBName string `json:"dbname"`
	} `json:"mongodb"`
}

var Conf Config

func Load() {
	arquivo, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Arquivo de configuração (config.json) não pode ser lido")
		os.Exit(0)
	}

	err = json.Unmarshal(arquivo, &Conf)
	if err != nil {
		fmt.Println("erro: ", err)
	}
}
