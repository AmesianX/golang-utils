package bot

import (
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// TELEGRAM :
const TELEGRAM = "telegram"

// Telegram :
type Telegram struct {
	Token   string
	ChatID  int64
	Channel string

	PostChanChan chan chan *Post
	Done         chan int

	botAPI        *tgbotapi.BotAPI
	config        tgbotapi.UpdateConfig
	updateChannel tgbotapi.UpdatesChannel
}

// NewTelegram :
func NewTelegram(token string, chatID int64, channel string) *Telegram {
	return &Telegram{Token: token, ChatID: chatID, Channel: channel, PostChanChan: make(chan chan *Post, 1), Done: make(chan int, 1)}
}

// Login :
func (t *Telegram) Login() error {
	t.config = tgbotapi.NewUpdate(0)
	t.config.Timeout = 60

	bot, e := tgbotapi.NewBotAPI(t.Token)
	if e != nil {
		return e
	}
	t.botAPI = bot

	t.updateChannel, e = t.botAPI.GetUpdatesChan(t.config)
	if e != nil {
		return e
	}

	return nil
}

// GetPostChanChan :
func (t *Telegram) GetPostChanChan() chan chan *Post {
	return t.PostChanChan
}

// Start :
func (t *Telegram) Start() {
	go func() {
		postChan := <-t.PostChanChan
		for {
			select {
			case req := <-t.updateChannel:
				postChan <- NewPost(TELEGRAM, t.Channel, req.Message.Text)
			case <-t.Done:
				break
			}
		}
	}()
}

// Send :
func (t Telegram) Send(post *Post) error {
	if strings.Compare(TELEGRAM, post.Messenger) == 0 {
		if strings.Compare(t.Channel, post.Channel) == 0 {
			telegramPost := tgbotapi.NewMessage(t.ChatID, post.Message)

			_, e := t.botAPI.Send(telegramPost)
			if e != nil {
				return e
			}
		}
	}

	return nil
}

// Shutdown :
func (t Telegram) Shutdown() {
	t.Done <- 1
}
