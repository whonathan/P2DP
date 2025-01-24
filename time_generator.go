package main

import (
	"math/rand"
	"sync"
	"time"
)

type TimeGenerator struct {
	mu       sync.Mutex
	lastTime time.Time
	rng      *rand.Rand
}

func NewTimeGenerator(seed int64) *TimeGenerator {
	return &TimeGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

func (g *TimeGenerator) NextTime() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now()

	if g.lastTime.IsZero() {
		startHour := 7 + g.rng.Intn(2)
		g.lastTime = time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			startHour,
			g.rng.Intn(60),
			g.rng.Intn(60),
			0,
			time.Local,
		)
		return g.lastTime.Format("01/02/2006 15:04:05")
	}

	nextGap := time.Duration(120+g.rng.Intn(120)) * time.Second
	g.lastTime = g.lastTime.Add(nextGap)
	return g.lastTime.Format("01/02/2006 15:04:05")
}
