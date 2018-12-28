package util

import (
	"errors"
	"time"
)

// Alarm :
type Alarm struct {
	Enable bool `json:"enable"`

	Hour int `json:"hour"`
	Min  int `json:"min"`
	Sec  int `json:"sec"`

	TimeChan chan time.Time `json:"-"`

	ticker *time.Ticker
}

// NewAlarm :
func NewAlarm(hour int, min int, sec int) *Alarm {
	return &Alarm{Enable: true, Hour: hour, Min: min, Sec: sec, TimeChan: make(chan time.Time, 1)}
}

// IsValid :
func (a *Alarm) IsValid() error {
	if a.Hour < 0 || a.Hour > 23 {
		return errors.New("Hour is invalid")
	}
	if a.Min < 0 || a.Min > 59 {
		return errors.New("Min is invalid")
	}
	if a.Sec < 0 || a.Sec > 59 {
		return errors.New("Sec is invalid")
	}
	return nil
}

// Start :
func (a *Alarm) Start() {
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), a.Hour, a.Min, a.Sec, 0, now.Location())

	if !t.After(now) {
		t = time.Date(now.Year(), now.Month(), now.Day()+1, a.Hour, a.Min, a.Sec, 0, now.Location())
	}

	a.ticker = time.NewTicker(t.Sub(now))

	go func() {
		t := <-a.ticker.C
		a.ticker.Stop()

		a.ticker = time.NewTicker(24 * time.Hour)

		a.TimeChan <- t

		for {
			select {
			case t := <-a.ticker.C:
				a.TimeChan <- t
			}
		}
	}()
}

// Stop :
func (a *Alarm) Stop() {
	if a != nil {
		if a.ticker != nil {
			a.ticker.Stop()
		}
	}
}
