package bot

import (
	"errors"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// TELEGRAM :
const TELEGRAM = "telegram"

// Telegram :
type Telegram struct {
	token   string
	chatID  int64
	channel string

	postChanChan chan chan *Post
	done         chan int

	botAPI        *tgbotapi.BotAPI
	config        tgbotapi.UpdateConfig
	updateChannel tgbotapi.UpdatesChannel
}

// NewTelegram :
func NewTelegram(token string, chatID int64, channel string) *Telegram {
	return &Telegram{token: token, chatID: chatID, channel: channel, postChanChan: make(chan chan *Post, 1), done: make(chan int, 1)}
}

// IsValid :
func (t *Telegram) IsValid() error {
	if len(t.token) == 0 {
		return errors.New("token is nil")
	}
	if t.chatID == 0 {
		return errors.New("chatID is 0")
	}

	return nil
}

// GetPostChanChan :
func (t *Telegram) GetPostChanChan() chan chan *Post {
	return t.postChanChan
}

// Login :
func (t *Telegram) Login() error {
	t.config = tgbotapi.NewUpdate(0)
	t.config.Timeout = 60

	bot, e := tgbotapi.NewBotAPI(t.token)
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

// Start :
func (t *Telegram) Start() {
	go func() {
		postChan := <-t.postChanChan
		for {
			select {
			case req := <-t.updateChannel:
				postChan <- NewPost(TELEGRAM, t.channel, req.Message.Text)
			case <-t.done:
				break
			}
		}
	}()
}

// Send :
func (t Telegram) Send(post *Post) error {
	if strings.Compare(TELEGRAM, post.Messenger) == 0 {
		if strings.Compare(t.channel, post.Channel) == 0 {
			telegramPost := tgbotapi.NewMessage(t.chatID, post.Message)

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
	t.done <- 1
}
