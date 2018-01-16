package main

import (
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

func DumpGoroutinesOnSigQUIT() {
	var debugChan = make(chan os.Signal, 1)
	signal.Notify(debugChan, syscall.SIGQUIT)
	go func() {
		select {
		case <-debugChan:
			pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
		}
	}()
}
