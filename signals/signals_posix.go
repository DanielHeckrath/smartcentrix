// +build !windows

package signals

import (
	"os"
	"os/signal"
	"syscall"
)

func Handle(quit QuitFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)

	for sig := range signals {
		switch sig {
		case syscall.SIGINT, syscall.SIGUSR1, syscall.SIGTERM:
			quit()
		}
	}
}
