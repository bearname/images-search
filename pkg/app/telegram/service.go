package telegram

import (
	"github.com/col3name/images-search/pkg/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
