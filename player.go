package main

import (
	"time"

	"github.com/lusingander/go-gif-viewer/image"
)

type player struct {
	delayMilliSec int // TODO: fix

	observers []func()
	stopPlay  chan bool
}

func newPlayer(i *image.GIFImage) *player {
	return &player{
		delayMilliSec: i.DelayMilliSec()[0], // TODO: fix
		observers:     make([]func(), 0),
	}
}

func (p *player) addObserver(f func()) {
	p.observers = append(p.observers, f)
}

func (p *player) play() {
	if p.stopPlay != nil {
		return
	}
	p.stopPlay = make(chan bool)
	go func() {
		t := time.NewTicker(time.Duration(p.delayMilliSec) * time.Millisecond)
		for {
			select {
			case <-t.C:
				for _, o := range p.observers {
					o()
				}
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
