package util

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func ExecCommand(process string, args ...string) error {
	call := exec.Command(process, args...)
	call.Stderr = os.Stderr
	call.Stdout = os.Stdout
	call.Stdin = os.Stdin

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-sigs:
			case <-done:
				break
			}
		}
	}()
	defer close(done)

	if err := call.Run(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
