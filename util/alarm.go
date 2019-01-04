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

	timeChan chan time.Time

	ticker *time.Ticker
}

// NewAlarm :
func NewAlarm(hour int, min int, sec int) *Alarm {
	return &Alarm{Enable: true, Hour: hour, Min: min, Sec: sec, timeChan: make(chan time.Time, 1)}
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

// GetTimeChan :
func (a *Alarm) GetTimeChan() chan time.Time {
	return a.timeChan
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

		a.timeChan <- t

		for {
			select {
			case t := <-a.ticker.C:
				a.timeChan <- t
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
