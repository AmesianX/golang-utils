package bot

import (
	"gopkg.in/telegram-bot-api.v4"
)

// TELEGRAM :
const TELEGRAM = "telegram"

// Telegram :
type Telegram struct {
	Enable bool

	Token  string
	ChatID int64

	PostChan chan chan *Post
	Done     chan int

	botAPI        *tgbotapi.BotAPI
	config        tgbotapi.UpdateConfig
	updateChannel tgbotapi.UpdatesChannel
}

// NewTelegram :
func NewTelegram(enable bool, token string, chatID int64) *Telegram {
	return &Telegram{Enable: enable, Token: token, ChatID: chatID, PostChan: make(chan chan *Post, 1), Done: make(chan int, 1)}
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

// Start :
func (t *Telegram) Start() {
	go func() {
		postChan := <-t.PostChan
		for {
			select {
			case req := <-t.updateChannel:
				postChan <- NewPost(TELEGRAM, req.Message.Text)
			case <-t.Done:
				break
			}
		}
	}()
}

// Send :
func (t Telegram) Send(post *Post) error {
	telegramPost := tgbotapi.NewMessage(t.ChatID, post.Message)

	_, e := t.botAPI.Send(telegramPost)
	if e != nil {
		return e
	}

	return nil
}

// Shutdown :
func (t Telegram) Shutdown() {
	t.Done <- 1
}
