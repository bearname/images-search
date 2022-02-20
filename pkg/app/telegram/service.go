package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"photofinish/pkg/domain"
)

type Notifier struct {
	api           tgbotapi.BotAPI
	supportChatId int64
}

func NewNotifier(api *tgbotapi.BotAPI, supportChatId int) *Notifier {
	n := new(Notifier)
	n.api = *api
	n.supportChatId = int64(supportChatId)
	return n
}

func (s *Notifier) Notify(msg domain.Message) error {
	message := tgbotapi.NewMessage(s.supportChatId, msg.Message)
	_, err := s.api.Send(message)
	return err
}
