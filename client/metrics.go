package main

import (
	"fmt"
	"sync"
	"time"
)

type Metrics struct {
	sync.Mutex
	durations      map[string]time.Duration
	durationsCount map[string]int64
	bytes          map[string]int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		durations:      make(map[string]time.Duration),
		bytes:          make(map[string]int64),
		durationsCount: make(map[string]int64),
	}
}

func (t *Metrics) AddDuration(key string, duration time.Duration) {
	t.Lock()
	defer t.Unlock()

	t.durations[key] = t.durations[key] + duration
	t.durationsCount[key] = t.durationsCount[key] + 1
}

func (t *Metrics) AddBytes(key string, bytes int64) {
	t.Lock()
	defer t.Unlock()

	t.bytes[key] = t.bytes[key] + bytes
}

func (t *Metrics) Reset() {
	t.Lock()
	defer t.Unlock()

	fmt.Printf("Average Durations: ")
	for key, value := range t.durations {
		amount := value / time.Duration(t.durationsCount[key])
		fmt.Printf("%s=%s  ", key, amount)
	}
	fmt.Printf("\n")

	fmt.Printf("Total Kilobytes: ")
	for key, value := range t.bytes {
		fmt.Printf("%s=%.2f ", key, (float64(value) / 1024.0))
	}
	fmt.Printf("\n")

	t.durations = make(map[string]time.Duration)
	t.bytes = make(map[string]int64)
	t.durationsCount = make(map[string]int64)
}
