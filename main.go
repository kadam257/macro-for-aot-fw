package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/go-vgo/robotgo"
)

var (
	shiftHeld       bool
	inShiftLock     bool = false // Starting state: out of shift lock
	targetEnemies   bool = false // Starting state: no target enemies
	macrosEnabled   bool = false // Starting state: macros disabled
	stateMutex      sync.Mutex
)

func main() {
	fmt.Println("Macro program started...")
	fmt.Println("Initial state: Out of shift lock, no target enemies")
	fmt.Println("Macros: DISABLED (press '=' to toggle)")
	fmt.Println("")
	fmt.Println("Controls:")
	fmt.Println("  '=' - Toggle macros ON/OFF")
	fmt.Println("  '[' - Cycle shift lock modes (no targeting → targeting → no targeting)")
	fmt.Println("  'l' - Exit shift lock mode")
	fmt.Println("  'Ctrl+C' - Exit")
	fmt.Println("")
	fmt.Println("Press '=' to enable macros when you're ready to use them!")

	keyboard := "/dev/input/event4"
	mouse := "/dev/input/event11"

	k, err := keylogger.New(keyboard)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	m, err := keylogger.New(mouse)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go keyListener(k)
	go keyListener(m)

	<-c
	fmt.Println("\nExiting...")
}

func shouldProcessMacros() bool {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	return macrosEnabled
}

func keyListener(k *keylogger.KeyLogger) {
	events := k.Read()

	for e := range events {
		switch e.Type {
		case keylogger.EvKey:
			keyStr := e.KeyString()
			
			if e.KeyPress() {
				switch keyStr {
				case "=":
					handleToggleKey()
				case "[":
					if shouldProcessMacros() {
						handleBracketKey()
					}
				case "L":
					if shouldProcessMacros() {
						handleLKey()
					}
				case "L_SHIFT", "R_SHIFT":
					if shouldProcessMacros() {
						handleShiftKeyDown()
					}
				}
			} else if e.KeyRelease() {
				switch keyStr {
				case "L_SHIFT", "R_SHIFT":
					if shouldProcessMacros() {
						handleShiftKeyUp()
					}
				}
			}
		}
	}
}

func handleBracketKey() {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	
	if !inShiftLock {
		// Not in shift lock: enter shift lock (no targeting)
		fmt.Println("[ pressed - entering shift lock mode")
		robotgo.KeyTap("shift")
		inShiftLock = true
		targetEnemies = false
		shiftHeld = false
	} else if inShiftLock && !targetEnemies {
		// In shift lock, no targeting: tap shift + hold shift (double shift to stay in shift lock)
		fmt.Println("[ pressed - enabling target enemies in shift lock (double shift)")
		robotgo.KeyTap("shift")
		time.Sleep(1 * time.Millisecond)
		robotgo.KeyDown("shift")
		targetEnemies = true
		shiftHeld = true
		// inShiftLock stays true
	} else if inShiftLock && targetEnemies {
		// In shift lock with targeting: disable targeting (back to shift lock only)
		fmt.Println("[ pressed - disabling target enemies, staying in shift lock")
		robotgo.KeyUp("shift")
		targetEnemies = false
		shiftHeld = false
	}
}


func handleLKey() {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	
	if inShiftLock {
		// Exit shift lock mode entirely
		if shiftHeld {
			robotgo.KeyUp("shift")
			shiftHeld = false
		}
		robotgo.KeyTap("shift") 
		
		inShiftLock = false
		targetEnemies = false
		
		fmt.Println("l pressed - exited shift lock mode")
	}
}

func handleShiftKeyDown() {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	
	if !shiftHeld {
		inShiftLock = !inShiftLock
		fmt.Printf("Manual shift pressed - now in %s mode\n", 
			map[bool]string{true: "shift lock", false: "normal"}[inShiftLock])
	}
}

func handleShiftKeyUp() {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	
	if !shiftHeld {
		fmt.Println("Manual shift released")
	}
}

func handleToggleKey() {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	
	macrosEnabled = !macrosEnabled
	status := map[bool]string{true: "ENABLED", false: "DISABLED"}[macrosEnabled]
	fmt.Printf("Macros: %s\n", status)
}
