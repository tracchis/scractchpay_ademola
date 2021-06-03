package process

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"syscall"
)

var sigHandler struct {
	sync.Once
	ctx  context.Context
	ch   chan os.Signal
	edge int64
}

// Context returns a process-level `context.Context`
// that is cancelled when the program receives a termination signal.
//
// A process should start a graceful shutdown process once this context is cancelled.
//
// This context should not be used as a parent to any requests,
// otherwise those requests will also be cancelled
// instead of being allowed to complete their work.
func Context() context.Context {
	sigHandler.Do(func() {
		sigHandler.ch = make(chan os.Signal, 1)
		sigHandler.ctx = signalHandler(context.Background())
	})

	return sigHandler.ctx
}

// signalHandler starts a long-lived goroutine, which will run in the background until the process ends.
// It returns a `context.Context` that will be canceled in the event of a signal being received.
// It will then continue to listen for more signals,
// If it receives three signals, then we presume that we are not shutting down properly, and panic with all stacktraces.
func signalHandler(parent context.Context) context.Context {
	ctx, cancel := context.WithCancel(parent)

	go func() {
		killChan := make(chan struct{}, 3)

		// `sigHandler.ch` should be setup before calling this function.
		// I want to avoid mutating the `sigHandler` in this function,
		// even though it should be safe anyways.
		signal.Notify(sigHandler.ch, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		for sig := range sigHandler.ch {
			fmt.Fprintln(os.Stderr, "received signal:", sig)

			switch sig {
			case syscall.SIGQUIT:
				debug.SetTraceback("all")
				panic("SIGQUIT")
			}

			// Prospectively set this value to 1, so we can maybe avoid injecting an unnecessary SIGTERM.
			atomic.StoreInt64(&sigHandler.edge, 1)

			cancel()

			select {
			case killChan <- struct{}{}: // killChan is not full, keep handling signals.
			default:
				// We have gotten three signals and we are still kicking,
				// panic and dump all stack traces.
				debug.SetTraceback("all")
				panic("not responding to signals")
			}
		}

		debug.SetTraceback("all")
		panic("signal handler channel unexpectedly closed")
	}()

	return ctx
}

// Shutdown starts any graceful shutdown processes waiting for `process.Context()` to be cancelled.
//
// Shutdown works by injecting a `syscall.SIGTERM` directly to the signal handler,
// which will cancel the `process.Context()` the same as a real SIGTERM.
//
// Shutdown returns an error only if it is unable to inject the signal.
// Notably, it is not an error to call Shutdown if we have already triggered a graceful shutdown.
//
// Shutdown does not wait for anything to finish before returning.
func Shutdown() error {
	// This will only return true if the value was 0 before setting it to 1.
	if !atomic.CompareAndSwapInt64(&sigHandler.edge, 0, 1) {
		// We have already triggered a graceful shutdown, so nothing to do.
		return nil
	}

	select {
	case sigHandler.ch <- syscall.SIGTERM:
		return nil
	default:
	}

	return errors.New("could not send signal")
}

// Quit ends the program as soon as possible, dumping a stacktrace of all goroutines.
//
// Quit works by injecting a `syscall.SIGQUIT` directly to the signal handler,
// which will cause a panic, and stacktrace of all goroutines the same as a real SIGQUIT.
//
// If Quit cannot inject the signal,
// it will setup an unrecoverable panic to occur.
//
// In all cases, Quit will not return.
func Quit() {
	select {
	case sigHandler.ch <- syscall.SIGQUIT:
	default:
		// We start up a separate goroutine for this to ensure that no `recover()` can block this panic.
		go func() {
			debug.SetTraceback("all")
			panic("process was force quit")
		}()
	}

	select {} // Block forever so we never return.
}
