package bot

// Bot :
type Bot interface {
	IsValid() error

	GetPostChanChan() chan chan *Post

	Login() error
	Start()
	Send(*Post) error
	Shutdown()
}
