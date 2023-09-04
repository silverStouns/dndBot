package request

import (
	"dndBot/internal/user_token"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type Request struct {
	Method    string
	Params    map[string]string
	Step      int
	UserToken string
	BotToken  string
	SentCode  bool
}

const (
	MetaTgID    string = "tgID"
	MetaChatID  string = "chatID"
	MetaMsgText string = "msgText"
)

func NewRequest(step int, method string, update tgbotapi.Update) *Request {
	m := make(map[string]string)
	return (&Request{
		Step:      step,
		Method:    method,
		UserToken: user_token.PrepareUserToken(update),
		Params:    m,
	}).
		SetMeta(MetaTgID, update.Message.From.ID).
		SetMeta(MetaChatID, update.Message.Chat.ID).
		SetMeta(MetaMsgText, update.Message.Text)
}

func (r *Request) SetMeta(key string, value interface{}) *Request {
	r.Params[key] = fmt.Sprintf("%v", value)

	return r
}

func (r *Request) SetParam(key string, value interface{}) *Request {
	r.Params[key] = fmt.Sprintf("%v", value)
	return r
}

func (r *Request) GetMeta(key string) (value string) {
	return r.Params[key]
}

func (r *Request) GetTgID() int64 {
	id, _ := strconv.ParseInt(r.GetMeta(MetaTgID), 10, 64)
	return id
}

func (r *Request) GetChatID() int64 {
	id, _ := strconv.ParseInt(r.GetMeta(MetaChatID), 10, 64)
	return id
}
