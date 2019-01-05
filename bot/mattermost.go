package bot

import (
	"errors"
	"net/url"
	"strings"

	"github.com/mattermost/mattermost-server/model"
)

// MATTERMOST :
const MATTERMOST = "mattermost"

// Mattermost :
type Mattermost struct {
	url      string
	username string
	password string
	team     string
	channel  string

	postChanChan chan chan *Post
	done         chan int

	client          *model.Client4
	botUser         *model.User
	botChannel      *model.Channel
	webSocketClient *model.WebSocketClient
}

// NewMattermost :
func NewMattermost(url string, username string, password string, team string, channel string) *Mattermost {
	return &Mattermost{url: url, username: username, password: password, team: team, channel: channel, postChanChan: make(chan chan *Post, 1), done: make(chan int, 1)}
}

// IsValid :
func (m *Mattermost) IsValid() error {
	if len(m.url) == 0 {
		return errors.New("url is nil")
	}

	_, e := url.Parse(m.url)
	if e != nil {
		return e
	}

	if len(m.username) == 0 {
		return errors.New("username is nil")
	}

	if len(m.password) == 0 {
		return errors.New("password is nil")
	}

	if len(m.team) == 0 {
		return errors.New("team is nil")
	}

	if len(m.channel) == 0 {
		return errors.New("channel is nil")
	}

	return nil
}

// GetPostChanChan :
func (m *Mattermost) GetPostChanChan() chan chan *Post {
	return m.postChanChan
}

// Login :
func (m *Mattermost) Login() error {
	m.client = model.NewAPIv4Client(m.url)

	if _, resp := m.client.GetOldClientConfig(""); resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	user, resp := m.client.Login(m.username, m.password)
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}
	m.botUser = user

	team, resp := m.client.GetTeamByName(m.team, "")
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	channel, resp := m.client.GetChannelByName(m.channel, team.Id, "")
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}
	m.botChannel = channel

	u, _ := url.Parse(m.url)

	webSocketClient, e := model.NewWebSocketClient4("wss://"+u.Hostname(), m.client.AuthToken)
	if e != nil {
		return errors.New(e.Message)
	}
	m.webSocketClient = webSocketClient

	m.webSocketClient.Listen()

	return nil
}

// Start :
func (m *Mattermost) Start() {
	go func() {
		postChan := <-m.postChanChan
		for {
			select {
			case eventChannel := <-m.webSocketClient.EventChannel:
				if eventChannel.Broadcast.ChannelId != m.botChannel.Id {
					continue
				}
				if eventChannel.Event != model.WEBSOCKET_EVENT_POSTED {
					continue
				}
				req := model.PostFromJson(strings.NewReader(eventChannel.Data["post"].(string)))
				if req != nil {
					if req.UserId == m.botUser.Id {
						continue
					}
				}

				postChan <- NewPost(MATTERMOST, m.channel, req.Message)
			case <-m.done:
				break
			}
		}
	}()
}

// Send :
func (m Mattermost) Send(post *Post) error {
	if strings.Compare(MATTERMOST, post.Messenger) == 0 {
		if strings.Compare(m.channel, post.Channel) == 0 {
			mattermostPost := &model.Post{}
			mattermostPost.ChannelId = m.botChannel.Id
			mattermostPost.Message = post.Message

			if _, resp := m.client.CreatePost(mattermostPost); resp.Error != nil {
				return errors.New(resp.Error.Message)
			}
		}
	}

	return nil
}

// Shutdown :
func (m Mattermost) Shutdown() {
	m.done <- 1
}
