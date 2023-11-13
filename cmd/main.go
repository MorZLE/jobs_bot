package main

import (
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/repository"
	"github.com/MorZLE/jobs_bot/service"
	"log"
)

func main() {
	cnf := config.NewConfig()

	st, err := repository.NewRepository(cnf)
	if err != nil {
		log.Fatal(err)
	}
	service.NewService(st)

	h := service.NewHandler(st, cnf)
	h.Start()
}
