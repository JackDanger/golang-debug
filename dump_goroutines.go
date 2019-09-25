package debug

import (
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

// usage:
//   kill -27 $pid_of_this_program
// or
//   kill -SIGPROF $pid_of_this_program
func init() {
	DumpGoroutinesOnSig(syscall.SIGPROF)
}

// DumpGoroutines prints the state of all goroutines to stdout
func DumpGoroutines() {
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}

// DumpGoroutinesOnSig registers a signal handler for the SIGPROF ('25', on most
// architectures) signal
func DumpGoroutinesOnSig(sig os.Signal) {
	var debugChan = make(chan os.Signal, 1)
	signal.Notify(debugChan, sig)
	go func() {
		for {
			select {
			case <-debugChan:
				DumpGoroutines()
			}
		}
	}()
}
