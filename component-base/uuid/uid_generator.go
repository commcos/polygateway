package uuid

import (
	"fmt"
	"time"
)

type AutoInc struct {
	start   int
	end     int
	queue   chan int
	running bool
}

var autoInc *AutoInc

func init() {
	autoInc = &AutoInc{
		start:   1,
		end:     10000,
		queue:   make(chan int, 10),
		running: true,
	}
	go autoInc.process()
}

func (ai *AutoInc) process() {
	for i := ai.start; ai.running; i = i + 1 {
		if i >= ai.end {
			i = ai.start
		}
		ai.queue <- i
	}
}

//GeneratesUID generates a new unique ID
func GeneratesUID(prefix, suffix string) string {
	currID := <-autoInc.queue
	if len(suffix) > 0 {
		return fmt.Sprintf("%s-%d.%d.%s", prefix, time.Now().Unix(), currID, suffix)
	}
	return fmt.Sprintf("%s-%d.%d", prefix, time.Now().Unix(), currID)
}

//Close close chan
func (ai *AutoInc) Close() {
	ai.running = false
	close(ai.queue)
}
