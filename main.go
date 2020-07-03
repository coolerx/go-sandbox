// No copyright

// sandbox of testing confusing golang features
package main

import (
	"fmt"
	"sync"
	"time"
)

const iKnowAboutPanic = true

func main() {
	selectSendOnClosed()
	goroutineSinglePanic()
}

// shows two usages
// panic when trying to send on closed channel
// panic recover by switch
func selectSendOnClosed() {
	defer func() {
		switch p := recover(); p {
		case nil:
			fmt.Println("selectSendOnClosed ok")
		default:
			fmt.Printf("selectSendOnClosed no: panic: %v\n", p)
		}
	}()

	ch := make(chan struct{})
	var choice string
	select {
	case ch <- struct{}{}:
		choice = "sent"
	default:
		choice = "N/A"
	}
	fmt.Println(choice)

	close(ch)

	select {
	case ch <- struct{}{}:
		choice = "sent"
	default:
		choice = "N/A"
	}
	fmt.Println(choice)
}

// single goroutine panic panics crashes whole program
func goroutineSinglePanic() {
	if !iKnowAboutPanic {
		defer func() {
			if p := recover(); p == "mustPanic" {
				fmt.Println("parent routine captures panic")
			} else {
				fmt.Println("parent routine can't capture panic")
			}
		}()
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go mustPanic(&wg)
	wg.Wait()

	wg.Add(1)
	go spinner(&wg, time.Second)
	wg.Wait()

	fmt.Println("\rreachable")
}

// simply panics
func mustPanic(wg *sync.WaitGroup) {
	if iKnowAboutPanic {
		defer func() {
			if p := recover(); p != nil {
				fmt.Println("recovered")
			}
		}()
	}

	defer wg.Done()

	panic("mustPanic")
}

// shows string to runes
// shows running fixed time
func spinner(wg *sync.WaitGroup, duration time.Duration) {
	defer wg.Done()

	runes := []rune(`\|/-`)
	doneCh := time.After(duration)

loop:
	for i := 0; ; i++ {
		select {
		case <-doneCh:
			break loop
		default:
			fmt.Printf("\r%c", runes[i%len(runes)])
			time.Sleep(time.Millisecond * 50)
		}
	}
}
