package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/MorZLE/jobs_bot/logger"
	"os"
	"strconv"
)

func NewConfig() *Config {
	return ParseFlags(&Config{})
}

type Config struct {
	DB       string
	BotToken string
	Dir      string
	Admin    int64
}

func ParseFlags(p *Config) *Config {
	var err error
	flag.StringVar(&p.DB, "a", "", "address db")
	flag.StringVar(&p.BotToken, "b", "", "BotToken")
	flag.Int64Var(&p.Admin, "c", 0, "id admin in telegram")
	flag.Parse()

	if DBAddr := os.Getenv("DB_ADDR"); DBAddr != "" {
		p.DB = DBAddr
	}

	if botToken := os.Getenv("BOT_TOKEN"); botToken != "" {
		p.BotToken = botToken
	}

	if admin := os.Getenv("FATHER_ADMIN"); admin != "" {
		num, err := strconv.ParseInt(admin, 10, 64)
		if err != nil {
			logger.Fatal("Не верный формат id в переменной окружения", errors.New("ErrFatherAdmin"))
		}
		p.Admin = num
	}
	if p.Admin == 0 {
		logger.Fatal("Не указан главный администратор", errors.New("ErrFatherAdmin"))
	}

	p.Dir, err = os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	return p
}
