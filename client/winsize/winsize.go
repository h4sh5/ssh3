//go:build !windows

package winsize

import (
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

func GetWinsize(tty *os.File) (ws WindowSize, err error) {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(tty.Fd()), uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)))
	if errno != 0 {
		err = errno
	}
	return ws, err
}

func MonWinsize(tty *os.File) <-chan WindowSize {
	// get a channel monitoring SIGWINCH
	ch := make(chan WindowSize, 1)
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGWINCH)
		defer signal.Stop(sigch)
		for range sigch {
			if ws, err := GetWinsize(tty); err == nil {
				ch <- ws
			}
		}
	}()
	return ch
}
