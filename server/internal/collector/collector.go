package collector

import (
	"fmt"
	"time"
)

type QueueEntity struct {
	id    uint
	title string
}

func RunCollector() {
	for {
		fmt.Println("Waiting for next collect...")
		time.Sleep(10 * time.Second)

	}
}
