package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"go.bug.st/serial.v1/enumerator"
)

var (
	ErrNotFound = errors.New("Port not found")
)

func check(vid string, pid string) (string, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return "", ErrNotFound
	}
	for _, port := range ports {
		if port.IsUSB {
			if port.VID == vid && port.PID == pid {
				return port.Name, nil
			}
		}
	}
	return "", ErrNotFound
}

var (
	vid = flag.String("vid", "", "USB vid to check for")
	pid = flag.String("pid", "", "USB pid to check for")
)

func wait() (string, error) {
	death := make(chan os.Signal, 1)
	port, err := check(*vid, *pid)
	if err == nil {
		return port, nil
	}

	signal.Notify(death, os.Interrupt, os.Kill)
	ticker := time.NewTicker(time.Second * 5) // check every 5 minutes
	fmt.Printf("Waiting for %s %s .", *vid, *pid)
	for {
		select {
		case <-death:
			os.Exit(0)
		case <-ticker.C:
			fmt.Print(".")
			port, err := check(*vid, *pid)
			if err != nil {
				continue
			}
			return port, nil
		}
	}
}

func main() {
	flag.Parse()
	port, err := wait()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	command := os.Args[6]
	args := os.Args[7:]
	for i, arg := range args {
		if arg == "{}" {
			args[i] = port
		}
	}
	if err := run(port, command, args); err != nil {
		panic(err)
	}
}