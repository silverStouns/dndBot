package config

import (
	"dndBot/internal/pkg/logger"
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

var (
	MainTgToken string
	MainBotLink string
)

// Команды призыва
const (
	Clean      = ""
	Start      = "/start"
	Summon     = "/Summon" // Призывает пул игроков(пул хранить в []string/file формате)
	Collection = "/Collection"
	Help       = "/help"
	ModulInfo  = "/moduls"
	GetChar    = "/get_char"
	Createchar = "/create_char"
)

const (
	CommandParametersSeparator = ":::"
)

func init() {
	configFile := flag.String("config", "internal/config/config.json", "this flag sets up the config file which will be used")
	flag.Parse()

	viper.SetConfigType("json")
	viper.SetConfigFile(*configFile)

	if err := viper.ReadInConfig(); err != nil {
		logger.Critical("Unable to read config file: %s", err)
		panic(err.Error())
	}
	if !viper.IsSet("TgToken") ||
		!viper.IsSet("BotLink") {
		panic(fmt.Sprintf("Wrong configs, no server authorization: used file %s", viper.ConfigFileUsed()))
	}
	MainTgToken = viper.GetString("TgToken")
	MainBotLink = viper.GetString("BotLink")
}
