package main

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/VacantionIT/coin-price-mon/internal/app/coinmonserver"
	"github.com/joho/godotenv"
)

var (
	configPath string
)

func init() {
	// run from cmd or root folder
	defConfigFile := "../../configs/coinmonserver.toml"
	if _, err := os.Stat(defConfigFile); err != nil {
		if os.IsNotExist(err) {
			defConfigFile = "configs/coinmonserver.toml"
		} else {
			log.Print(err)
		}
	}

	flag.StringVar(&configPath, "config-path", defConfigFile, "path to config")

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	flag.Parse()

	config := coinmonserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	s := coinmonserver.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
