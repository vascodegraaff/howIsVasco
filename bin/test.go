// package main

// import (
//     "fmt"
//     "github.com/looplab/fsm"
// )

// func main() {
//     fsm := fsm.NewFSM(
//         "start",
//         fsm.Events{
//             {Name: "1", Src: []string{"closed"}, Dst: "open"},
//             {Name: "close", Src: []string{"open"}, Dst: "closed"},
//         },
//         fsm.Callbacks{},
//     )

//     fmt.Println(fsm.Current())

//     err := fsm.Event("open")
//     if err != nil {
//         fmt.Println(err)
//     }

//     fmt.Println(fsm.Current())

//     err = fsm.Event("close")
//     if err != nil {
//         fmt.Println(err)
//     }

//     fmt.Println(fsm.Current())
// }
package main

import (
	"fmt"
	"time"
	"context"
)

func InputAsync() string {
	fmt.Println("Warming up ...")

	var data string
	fmt.Printf("input a word\n")
	fmt.Scan(&data)
	// time.Sleep(3 * time.Second)
	fmt.Println("Done ...")
	return data
}

// Future interface has the method signature for await
type Future interface {
	Await() interface{}
}

type future struct {
	await func(ctx context.Context) interface{}
}

func (f future) Await() interface{} {
	return f.await(context.Background())
}

// Exec executes the async function
func Exec(f func() interface{}) Future {
	var result interface{}
	c := make(chan struct{})
	go func() {
		defer close(c)
		result = f()
	}()
	return future{
		await: func(ctx context.Context) interface{} {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return result
			}
		},
	}
}

func getInput(input chan string) {
	for {
		var data string
		fmt.Println("input a string")
		fmt.Scan(&data)
		input <- data
	}
}

func main() {
	fmt.Println("Let's start ...")
	input := make(chan string, 1)
	go getInput(input)

	for {
		fmt.Println("input something")
		select {
		case i := <-input:
			fmt.Println("result")
			fmt.Println(i)

		case <-time.After(4000 * time.Millisecond):
			fmt.Println("timed out")
		}
	}
	// future := Exec(func() interface{} {
	// 	return InputAsync()
	// })
	// fmt.Println("Done is running ...")
	// val := future.Await()
	// fmt.Println(val)
}
