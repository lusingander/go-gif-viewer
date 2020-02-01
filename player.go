package main

import (
	"time"

	"github.com/lusingander/go-gif-viewer/image"
)

type player struct {
	delayMilliSec int // TODO: fix
	speed         float64

	observers   []func()
	stopPlay    chan bool
	changeSpeed chan float64
}

func newPlayer(i *image.GIFImage, speed float64) *player {
	return &player{
		delayMilliSec: i.DelayMilliSec()[0], // TODO: fix
		speed:         speed,
		observers:     make([]func(), 0),
		changeSpeed:   make(chan float64),
	}
}

func (p *player) addObserver(f func()) {
	p.observers = append(p.observers, f)
}

func (p *player) calcWaitDuration() time.Duration {
	ms := float64(p.delayMilliSec) * float64(time.Millisecond) / p.speed
	return time.Duration(ms)
}

func (p *player) play() {
	if p.stopPlay != nil {
		return
	}
	p.stopPlay = make(chan bool)
	go func() {
		t := time.NewTicker(p.calcWaitDuration())
		for {
			select {
			case <-t.C:
				for _, o := range p.observers {
					o()
				}
			case <-p.changeSpeed:
				t = time.NewTicker(p.calcWaitDuration())
			case <-p.stopPlay:
				t.Stop()
				return
			}
		}
	}()
}

func (p *player) pause() {
	if p.stopPlay == nil {
		return
	}
	p.stopPlay <- true
	p.stopPlay = nil
}

func (p *player) playing() bool {
	return p.stopPlay != nil
}

func (p *player) setSpeed(s float64) {
	p.speed = s
	if p.playing() {
		p.changeSpeed <- s
	}
}
