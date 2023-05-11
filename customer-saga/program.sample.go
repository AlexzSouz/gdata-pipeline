package main

import (
	"context"
	"fmt"
	"time"
)

func initMain() {
	timeout := 1500 * time.Microsecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ctx = context.WithValue(ctx, "label", "Value:")
	work := make(chan int)

	go processNumber(ctx, work)

	for i := 0; i < 3; i++ {
		work <- i
	}

	time.Sleep(1500 * time.Millisecond)
	fmt.Println("Work completed")
}

func processNumber(ctx context.Context, work chan int) {
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				fmt.Println("Error:", err)
			}
			fmt.Println("Context terminated")
			return
		case r := <-work:
			fmt.Println(ctx.Value("label"), r)
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
