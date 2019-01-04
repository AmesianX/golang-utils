package bot

// Post :
type Post struct {
	Messenger string
	Channel   string
	Message   string
}

// NewPost :
func NewPost(messenger string, channel string, message string) *Post {
	return &Post{Messenger: messenger, Channel: channel, Message: message}
}
