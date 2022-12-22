package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	totalInventory int
	inventorylock  sync.Mutex
	wg             sync.WaitGroup
)

func AddToInventory(ctx context.Context, inventory int) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	go func() {
		wg.Add(1)
		if err := process(inventory); err != nil {
			log.Print("Inventory is not valid")
		}
		wg.Done()
	}()
}

func process(inventory int) error {
	if inventory < 0 {
		return errors.New("inventory could not be negative")
	}

	// It takes two seconds to process the request.
	time.Sleep(time.Second * 2)

	log.Print("inventory proceed successfully")
	inventorylock.Lock()
	totalInventory += inventory
	inventorylock.Unlock()
	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	addInventories(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			cancel()
			wg.Wait()
			fmt.Println("canceled:", totalInventory)
			os.Exit(0)
		}
	}()

	wg.Wait()
	fmt.Println("finished:", totalInventory)
}

func addInventories(ctx context.Context) {
	for i := 0; i < 10; i++ {
		AddToInventory(ctx, i)
	}
}
