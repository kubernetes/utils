// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package inotify

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	child := exec.CommandContext(ctx, "sleep", "10")
	err := child.Start()
	if err != nil {
		t.Fatalf("exec sleep failed: %v", err)
	}
	go child.Wait()

	fdPath := filepath.Join("/proc", strconv.Itoa(child.Process.Pid), "fd")
	infos, err := ioutil.ReadDir(fdPath)
	if err != nil {
		t.Fatalf("read %s failed: %v", fdPath, err)
	}

	for _, info := range infos {
		path := filepath.Join(fdPath, info.Name())
		if info.Mode()&os.ModeSymlink != 0 {
			t.Logf("skipping %s", path)
		}
		target, err := os.Readlink(path)
		if err != nil {
			t.Errorf("readlink %s failed: %v", path, err)
		}
		if strings.Contains(target, "inotify") {
			t.Errorf("leaked inotify file descriptor: %s -> %s", path, target)
		}
	}
}
