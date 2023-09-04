package bot

import (
	"dndBot/internal/config"
	"dndBot/internal/database"
	"dndBot/internal/pkg/logger"
	"dndBot/internal/pkg/session"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
	"time"
)

var AllowedMethods = map[string]func(tgbotapi.Update) []tgbotapi.Chattable{
	"/help":                           HelpHandler,
	"/start":                          StartHandler,
	"/Summon":                         SummonHandler,
	"/Collection":                     CollectionHandler,
	"/help@DragoNs_EmperoR_bot":       HelpHandler,
	"/start@DragoNs_EmperoR_bot":      StartHandler,
	"/Summon@DragoNs_EmperoR_bot":     SummonHandler,
	"/Collection@DragoNs_EmperoR_bot": CollectionHandler,
	"/moduls":                         ModulInfoHandler,
	"/moduls@DragoNs_EmperoR_bot":     ModulInfoHandler,
	"/get_char":                       GetCharHandler,
	"/not_push":                       VideoHandler,
	"/not_push@DragoNs_EmperoR_bot":   VideoHandler,
	"/drop_dice":                      DropDiceHandler,
	"/drop_dice@DragoNs_EmperoR_bot":  DropDiceHandler,
	"/what":                           AudioHandler,
	"/what@DragoNs_EmperoR_bot":       AudioHandler,
	"/create_char":                    CreateCharHandler,
}
var Users map[int64]string
var timeActive time.Time

type User struct {
	Id        int64
	FirstName string
}

var DBconn *database.DBConnector

func ListenBot() {
	var err error
	DBconn, err = database.GetNewDBConnector()
	if err != nil {
		logger.Critical("server, err: %s", err.Error())
		os.Exit(1)
	}

	for {
		time.Sleep(time.Second)
		runBot()
	}
}

var PollAnswer map[string]string

func runBot() {
	bot, botInitError := tgbotapi.NewBotAPI(config.MainTgToken)
	if botInitError != nil {
		logger.Critical(botInitError.Error())
		return
	}
	bot.Debug = false
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализируем канал, куда будут прилетать обновления от API
	var ucfg = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates := bot.GetUpdatesChan(ucfg)
	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	for update := range updates {
		handleUpdate(bot, update)
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	defer OwnRecover()

	if update.PollAnswer != nil {
		timeActive = time.Now() // Запоминанием время первого ответа
		for _, answer := range update.PollAnswer.OptionIDs {
			if answer == 0 {
				Users[update.PollAnswer.User.ID] = update.PollAnswer.User.FirstName
			}
		}

	}
	if update.Message == nil && update.CallbackQuery == nil {
		return
	}

	if update.Message.From.IsBot {
		return
	}

	mergeCallbackIntoOrdinaryMessage(&update)

	arrayMessages := messageHandler(update)
	for _, response := range arrayMessages {
		sentMessage, err := bot.Send(response)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		if sentMessage.Document != nil {
			sentMessage.Text += "{{document}}"
		}
		if sentMessage.Photo != nil {
			sentMessage.Text += "{{image}}"
		}
		logger.Trace("[to %s]\n%s ", prepareUserNameForLogging(update), sentMessage.Text)
	}
}

func OwnRecover() {
	if r := recover(); r != nil {
		logger.Critical(fmt.Sprint(r))
	}
}

func messageHandler(update tgbotapi.Update) (messages []tgbotapi.Chattable) {
	request, continueHandlingRequest := session.GetSession(update)
	if update.Message.Photo != nil {
		return []tgbotapi.Chattable{tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста используйте встроенные команды")}
	}
	params := strings.Split(update.Message.Text, config.CommandParametersSeparator)
	if len(params) == 0 {
		return []tgbotapi.Chattable{}
	}

	if handler, foundHandler := AllowedMethods[params[0]]; foundHandler {
		if continueHandlingRequest {
			messages = append(messages, exitHandler(update)...)
		}
		return append(messages, handler(update)...)
	}

	if !continueHandlingRequest {
		logger.Trace("continueHandlingRequest")
		return append(messages, unknownMethodHandler(update)...)
	}
	handler, foundHandler := AllowedMethods[request.Method]
	if !foundHandler {
		logger.Trace("continueHandlingRequest")
		return []tgbotapi.Chattable{tgbotapi.NewMessage(update.Message.Chat.ID, "Непредвиденная ошибка")}
	}
	return append(messages, handler(update)...)
}

var unknownMethodHandler = func(update tgbotapi.Update) []tgbotapi.Chattable {
	return []tgbotapi.Chattable{}
}

var exitHandler = func(update tgbotapi.Update) []tgbotapi.Chattable {
	logger.Debug("Завершение сессии")
	session.FinishSession(update)
	return []tgbotapi.Chattable{} //nolint:gofmt
}

func mergeCallbackIntoOrdinaryMessage(update *tgbotapi.Update) {
	logger.Debug("calBack:%+v", update.Message)
	if update.CallbackQuery != nil {
		update.Message = update.CallbackQuery.Message
		update.Message.From = update.CallbackQuery.From
		update.Message.Text = update.CallbackQuery.Data
		update.Message.Text = update.CallbackQuery.InlineMessageID
	}
}

func prepareUserNameForLogging(update tgbotapi.Update) string {
	if update.Message.From.UserName != "" {
		return fmt.Sprintf("%s (id:%d)", update.Message.From.UserName, update.Message.From.ID)
	}
	return fmt.Sprintf("%s (id:%d)", "hidden", update.Message.From.ID)
}

func init() {
	Users = make(map[int64]string)
}

func Cleaner() {
	for {
		if time.Now().After(timeActive.AddDate(0, 0, 5)) {
			Users = make(map[int64]string)
		}
	}

}
