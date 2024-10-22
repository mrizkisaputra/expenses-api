package main

import (
	"github.com/mrizkisaputra/expenses-api/config"
	"log"
	"os"
)

func main() {
	_, err := config.NewAppConfig(os.Getenv("config"))
	if err != nil {
		log.Fatal(err)
	}

}
