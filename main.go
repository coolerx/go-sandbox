// No copyright

// sandbox of testing confusing golang features
package main

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/jszwec/csvutil"
)

const iKnowAboutPanic = true

func main() {
	selectSendOnClosed()
	goroutineSinglePanic()
	csvSandbox()
	deferReturn()
	mimicOverride()
	timeSandbox()
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
			fmt.Printf("selectSendOnClosed no: panic: %T %v\n", p, p)
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

func csvSandbox() {
	var csvInput = []byte(`
ID,name,age,CreatedAt
1,jacek,26,2012-04-01T15:00:00Z
2,john,,0001-01-01T00:00:00Z`)

	type UserID int64
	type User struct {
		ID   UserID
		Name string `csv:"name"`
		Age  int    `csv:"age,omitempty"`
		// Age       int `csv:"age"`
		CreatedAt time.Time
	}

	var users []User
	if err := csvutil.Unmarshal(csvInput, &users); err != nil {
		fmt.Println("error:", err)
	}

	for _, u := range users {
		fmt.Printf("%+v\n", u)
	}

	type UnitID int64

	var user1 UserID = 1
	var unit1 UnitID = 1
	if user1 == UserID(unit1) {
		fmt.Printf("%d %d\n", user1, unit1)
	}
}

func deferReturn() {
	finalErr := func() (err error) {
		defer func() {
			err = errors.New("deferred error")
		}()
		err = errors.New("original error")
		return
	}()
	fmt.Println("err:", finalErr)
}

type talker interface {
	talk() string
	talk2() string
}

type defaultTalker struct{}

func (t defaultTalker) talk() string {
	return "default"
}

func (t defaultTalker) talk2() string {
	return "default2"
}

type basicTalker struct{ defaultTalker }

type overrideTalker struct{ defaultTalker }

func (t overrideTalker) talk() string {
	return "override"
}

func mimicOverride() {
	bt := talker(basicTalker{})
	ot := talker(overrideTalker{})

	fmt.Println(reflect.TypeOf(bt), bt.talk())
	fmt.Println(reflect.TypeOf(ot), ot.talk())
	fmt.Println(reflect.TypeOf(bt), bt.talk2())
	fmt.Println(reflect.TypeOf(ot), ot.talk2())
}

func timeSandbox() {
	now := time.Now()

	fmt.Printf("now: %v\n", now)
	fmt.Println(now.Format(time.RFC3339))
	fmt.Println(now.Format(time.RFC850))

	fmt.Printf("unix now: %v\n", now.Unix())
	fmt.Printf("unix utc now: %v\n", now.UTC().Unix())
}
