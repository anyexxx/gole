package console

import (
	"bufio"
	"fmt"
	"me.qqtu.game/router"
	"me.qqtu.game/server"
	"os"
	"runtime"
	"syscall"
)

var (
	cpus         = "cpu num"
	routines     = "goroutines"
	usernum      = "user num"
	reloadrouter = "reload router"
	memusage     = "memory usage"
	quit         = "quit"
)

var (
	command = make([]byte, 1024)
	reader  = bufio.NewReader(os.Stdin)
)

func Console() {
	for {
		fmt.Print(router.CurrHost + ":" + router.CurrPort + ">>")
		command, _, _ = reader.ReadLine()
		switch string(command) {

		case quit:
			fmt.Println("Server stopped.")
			os.Exit(0)
		case cpus:
			fmt.Println("The number of CPUs currently in use: ", runtime.NumCPU())
		case routines:
			fmt.Println("Current number of goroutines: ", runtime.NumGoroutine())
		case usernum:
			fmt.Println("The number of clients currently online is ", server.NumberOfUser())
		case memusage:
			us := &syscall.Rusage{}
			err := syscall.Getrusage(syscall.RUSAGE_SELF, us)
			if err != nil {
				fmt.Println("Get usage error: ", err, "\n")
			} else {
				fmt.Println("Memory Usage: %f MB\n\n", float64(us.Maxrss)/1024/1024)
			}
		case reloadrouter:
			router.ReloadConfig()
			fmt.Println("Reload router config done!")
		}
	}
	defer fmt.Println("")
}
