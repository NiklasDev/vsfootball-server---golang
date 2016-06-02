package cron

import (
	"code.google.com/p/gcfg"
	"github.com/robfig/cron"
	"fmt"
	// "log"
	"time"
)

type Config struct {
	Game struct {
		Warning 	int
		ResignGame 	int
	}
}

func ReadConfig(file string) Config {
	var conf Config
	err := gcfg.ReadFileInto(&conf, "config.ini")
	if err != nil {
		fmt.Println("Failed to parse gcfg data: %s", err)

		conf.Game.Warning = 3
		conf.ResignGame = 7
	} else {
		fmt.Println("Game.Warning : %d \n Game.ResignGame %d", conf.Game.Warning, conf.Game.ResignGame)
	}

	return conf
}

func NotifyPlayers() {

} 

func DailyJob() {
	conf := ReadConfig('./config.ini')

	c := cron.New()
	c.AddJob("0 0 2 * * *", FuncJob(NotifyPlayers))
}

