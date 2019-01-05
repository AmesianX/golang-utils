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
	URL      string
	Username string
	Password string
	Team     string
	Channel  string

	PostChanChan chan chan *Post
	Done         chan int

	client          *model.Client4
	botUser         *model.User
	botChannel      *model.Channel
	webSocketClient *model.WebSocketClient
}

// NewMattermost :
func NewMattermost(url string, username string, password string, team string, channel string) *Mattermost {
	return &Mattermost{URL: url, Username: username, Password: password, Team: team, Channel: channel, PostChanChan: make(chan chan *Post, 1), Done: make(chan int, 1)}
}

// Login :
func (m *Mattermost) Login() error {
	m.client = model.NewAPIv4Client(m.URL)

	if _, resp := m.client.GetOldClientConfig(""); resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	user, resp := m.client.Login(m.Username, m.Password)
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}
	m.botUser = user

	team, resp := m.client.GetTeamByName(m.Team, "")
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	channel, resp := m.client.GetChannelByName(m.Channel, team.Id, "")
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}
	m.botChannel = channel

	u, _ := url.Parse(m.URL)

	webSocketClient, e := model.NewWebSocketClient4("wss://"+u.Hostname(), m.client.AuthToken)
	if e != nil {
		return errors.New(e.Message)
	}
	m.webSocketClient = webSocketClient

	m.webSocketClient.Listen()

	return nil
}

// GetPostChanChan :
func (m *Mattermost) GetPostChanChan() chan chan *Post {
	return m.PostChanChan
}

// Start :
func (m *Mattermost) Start() {
	go func() {
		postChan := <-m.PostChanChan
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

				postChan <- NewPost(MATTERMOST, m.Channel, req.Message)
			case <-m.Done:
				break
			}
		}
	}()
}

// Send :
func (m Mattermost) Send(post *Post) error {
	if strings.Compare(MATTERMOST, post.Messenger) == 0 {
		if strings.Compare(m.Channel, post.Channel) == 0 {
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
	m.Done <- 1
}
