package config

import (
	"flag"
	"fmt"
	"os"
)

func NewConfig() *Config {
	return ParseFlags(&Config{})
}

type Config struct {
	DB       string
	BotToken string
	Dir      string
}

func ParseFlags(p *Config) *Config {

	var err error
	flag.StringVar(&p.DB, "a", "", "address db")
	flag.StringVar(&p.BotToken, "b", "", "BotToken")
	flag.Parse()

	if DBAddr := os.Getenv("DB_ADDR"); DBAddr != "" {
		p.DB = DBAddr
	}

	if botToken := os.Getenv("BOT_TOKEN"); botToken != "" {
		p.BotToken = botToken
	}

	fmt.Println(p.BotToken)

	p.Dir, err = os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	return p
}
