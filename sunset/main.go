package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/davecheney/gpio"
	"github.com/davecheney/gpio/rpi"
)

// Show sunrise and sunset for first 5 days of June in LA
func main() {

	fmt.Println(GetSunset())
	fmt.Println("SUNSET!")

	// set GPIO25 to output mode
	pin, err := gpio.OpenPin(rpi.GPIO17, gpio.ModeOutput)
	if err != nil {
		fmt.Printf("Error opening pin! %s\n", err)
		return
	}

	// turn the led off on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Printf("\nClearing and unexporting the pin.\n")
			pin.Clear()
			pin.Close()
			os.Exit(0)
		}
	}()

	for {
		pin.Set()
		time.Sleep(1000 * time.Millisecond)
		pin.Clear()
		time.Sleep(1000 * time.Millisecond)
	}
}
