package main

import (
	"log"
	"sync/atomic"
	"time"
)

type Stats struct {
	serverStartTime    time.Time
	currentConnections uint32
	maxConnections     uint32
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) SetServerStartTime() {
	s.serverStartTime = time.Now()
}

func (s *Stats) IncreaseCurrentConnections() {
	atomic.AddUint32(&s.currentConnections, 1)

	currentCount := atomic.LoadUint32(&s.currentConnections)
	maxCount := atomic.LoadUint32(&s.maxConnections)

	if currentCount > maxCount {
		atomic.AddUint32(&s.maxConnections, currentCount-maxCount)
	}
}

func (s *Stats) DecreaseCurrentConnections() {
	if atomic.LoadUint32(&s.currentConnections) > 1 {
		// We use this hack because we cant use uint32(-1), unsinged int's cant go below 0
		decreaseHack := int32(-1)
		atomic.AddUint32(&s.currentConnections, uint32(decreaseHack))
	}
}

func (s *Stats) PrintStats() {
	log.Printf("[STATS]\t Server uptime: %s", time.Now().Sub(s.serverStartTime))
	log.Printf("[STATS]\t Current connections: %d", s.currentConnections)
	log.Printf("[STATS]\t Maximum connections: %d", s.maxConnections)
}
