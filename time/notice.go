package time

import (
	"errors"
	"time"
)

// Notice 알림
type Notice struct {
	Enable bool `json:"enable"`

	Hour int `json:"hour"`
	Min  int `json:"min"`
	Sec  int `json:"sec"`

	TimeChan chan time.Time

	ticker *time.Ticker
}

func NewNotice(hour int, min int, sec int) *Notice {
	return &Notice{Enable: true, Hour: hour, Min: min, Sec: sec, TimeChan: make(chan time.Time, 1)}
}

func (n *Notice) IsValid() error {
	if n.Hour < 0 || n.Hour > 23 {
		return errors.New("Hour is invalid")
	}
	if n.Min < 0 || n.Min > 59 {
		return errors.New("Min is invalid")
	}
	if n.Sec < 0 || n.Sec > 59 {
		return errors.New("Sec is invalid")
	}
	return nil
}

// Start 정해진 시간에 TimeChan을 통해 알림
func (n *Notice) Start() {
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), n.Hour, n.Min, n.Sec, 0, now.Location())

	if !t.After(now) {
		t = time.Date(now.Year(), now.Month(), now.Day()+1, n.Hour, n.Min, n.Sec, 0, now.Location())
	}

	n.ticker = time.NewTicker(t.Sub(now))

	go func() {
		t := <-n.ticker.C
		n.ticker.Stop()

		n.ticker = time.NewTicker(24 * time.Hour)

		n.TimeChan <- t

		for {
			select {
			case t := <-n.ticker.C:
				n.TimeChan <- t
			}
		}
	}()
}

// Stop 알림 중지
func (n *Notice) Stop() {
	if n != nil {
		if n.ticker != nil {
			n.ticker.Stop()
		}
	}
}
