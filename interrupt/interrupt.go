/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package interrupt

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// terminationSignals are the default signals that cause the program to exit in the
// supported platforms (linux, darwin, windows).
var terminationSignals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

// Handler guarantees execution of notifications after a critical section (the function passed
// to a Run method), even in the presence of process termination. It guarantees exactly once
// invocation of the provided notify functions.
type Handler struct {
	notify []func()
	final  func(os.Signal)
	once   sync.Once
}

// Chain creates a new handler that invokes all notify functions when the critical section exits
// and then invokes the optional handler's notifications. This allows critical sections to be
// nested without losing exactly once invocations. Notify functions can invoke any cleanup needed
// but should not exit (which is the responsibility of the parent handler).
func Chain(handler *Handler, notify ...func()) *Handler {
	if handler == nil {
		return New(nil, notify...)
	}
	return New(handler.signal, append(notify, handler.close)...)
}

// New creates a new handler that guarantees all notify functions are run after the critical
// section exits (or is interrupted by the OS), then invokes the final handler. Notify
// functions are invoked serially in the order passed to New method. If no final
// handler is specified, the default final is `os.Exit(1)`. A handler can only be used for
// one critical section. The final handler is given complete control over reacting to the
// signal - it may exit the process or ignore the signal as needed. The following signals
// are watched: SIGHUP, SIGINT, SIGTERM, and SIGQUIT.
func New(final func(os.Signal), notify ...func()) *Handler {
	return &Handler{
		final:  final,
		notify: notify,
	}
}

// close executes all the notification handlers if they have not yet been executed.
func (h *Handler) close() {
	h.once.Do(func() {
		for _, fn := range h.notify {
			fn()
		}
	})
}

// signal is called when an os.Signal is received, and guarantees that all notifications
// are executed, then the final handler is executed. This function should only be called once
// per Handler instance.
func (h *Handler) signal(s os.Signal) {
	h.once.Do(func() {
		for _, fn := range h.notify {
			fn()
		}
		if h.final == nil {
			os.Exit(1)
		}
		h.final(s)
	})
}

// Run ensures that any notifications are invoked after the provided fn exits (even if the
// process is interrupted by an OS termination signal). Notifications are only invoked once
// per Handler instance, so calling Run more than once will not behave as the user expects.
func (h *Handler) Run(fn func() error) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, terminationSignals...)
	defer func() {
		signal.Stop(ch)
		close(ch)
	}()
	go func() {
		sig, ok := <-ch
		if !ok {
			return
		}
		h.signal(sig)
	}()
	defer h.close()
	return fn()
}
