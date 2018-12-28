package bot

// Post :
type Post struct {
	Messenger string
	Message   string
}

// NewPost :
func NewPost(messenger string, message string) *Post {
	return &Post{Messenger: messenger, Message: message}
}
