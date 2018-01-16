package debug

import (
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

func DumpGoroutines() {
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}

func DumpGoroutinesOnSig(sig os.Signal) {
	var debugChan = make(chan os.Signal, 1)
	signal.Notify(debugChan, sig)
	go func() {
		select {
		case <-debugChan:
			DumpGoroutines()
		}
	}()
}
func DumpGoroutinesOnSigQUIT() {
	DumpGoroutinesOnSig(syscall.SIGQUIT)
}
