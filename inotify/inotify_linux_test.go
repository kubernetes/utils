// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package inotify

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

func TestInotifyEvents(t *testing.T) {
	// Create an inotify watcher instance and initialize it
	watcher, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %s", err)
	}

	dir, err := ioutil.TempDir("", "inotify")
	if err != nil {
		t.Fatalf("TempDir failed: %s", err)
	}
	defer os.RemoveAll(dir)

	// Add a watch for "_test"
	err = watcher.Watch(dir)
	if err != nil {
		t.Fatalf("Watch failed: %s", err)
	}

	// Receive errors on the error channel on a separate goroutine
	go func() {
		for err := range watcher.Error {
			t.Fatalf("error received: %s", err)
		}
	}()

	testFile := dir + "/TestInotifyEvents.testfile"

	// Receive events on the event channel on a separate goroutine
	eventstream := watcher.Event
	var eventsReceived int32
	done := make(chan bool)
	go func() {
		for event := range eventstream {
			// Only count relevant events
			if event.Name == testFile {
				atomic.AddInt32(&eventsReceived, 1)
				t.Logf("event received: %s", event)
			} else {
				t.Logf("unexpected event received: %s", event)
			}
		}
		done <- true
	}()

	// Create a file
	// This should add at least one event to the inotify event queue
	_, err = os.OpenFile(testFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatalf("creating test file: %s", err)
	}

	// We expect this event to be received almost immediately, but let's wait 1 s to be sure
	time.Sleep(1 * time.Second)
	if atomic.AddInt32(&eventsReceived, 0) == 0 {
		t.Fatal("inotify event hasn't been received after 1 second")
	}

	// Try closing the inotify instance
	t.Log("calling Close()")
	watcher.Close()
	t.Log("waiting for the event channel to become closed...")
	select {
	case <-done:
		t.Log("event channel closed")
	case <-time.After(1 * time.Second):
		t.Fatal("event stream was not closed after 1 second")
	}
}

func TestInotifyClose(t *testing.T) {
	watcher, _ := NewWatcher()
	watcher.Close()

	done := make(chan bool)
	go func() {
		watcher.Close()
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(50 * time.Millisecond):
		t.Fatal("double Close() test failed: second Close() call didn't return")
	}

	err := watcher.Watch(os.TempDir())
	if err == nil {
		t.Fatal("expected error on Watch() after Close(), got nil")
	}
}

func TestInotifyFdLeak(t *testing.T) {
	watcher, _ := NewWatcher()
	defer watcher.Close()

	child := exec.Command("sleep", "1")
	err := child.Start()
	if err != nil {
		t.Fatalf("exec sleep failed: %v", err)
	}

	pid := child.Process.Pid
	fds, err := ioutil.ReadDir(fmt.Sprintf("/proc/%d/fd", pid))
	_ = syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		t.Fatalf("read procfs of %d failed: %v", pid, err)
	}
	var actualFds []string
	for _, fd := range fds {
		actualFds = append(actualFds, fd.Name())
	}

	// stdin, stdout, stderr
	expectFds := []string{"0", "1", "2"}
	if !reflect.DeepEqual(expectFds, actualFds) {
		t.Fatalf("expected fds: %+v, actual fds %+v", expectFds, actualFds)
	}
}
