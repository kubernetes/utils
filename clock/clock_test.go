/*
Copyright 2026 The Kubernetes Authors.

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

package clock

import (
	"testing"
	"time"
)

func TestRealClockNow(t *testing.T) {
	rc := RealClock{}
	before := time.Now()
	got := rc.Now()
	after := time.Now()
	if got.Before(before) || got.After(after) {
		t.Fatalf("Now() = %v, want between %v and %v", got, before, after)
	}
}

func TestRealClockSince(t *testing.T) {
	rc := RealClock{}
	start := time.Now().Add(-100 * time.Millisecond)
	d := rc.Since(start)
	if d < 100*time.Millisecond {
		t.Fatalf("Since() = %v, want at least 100ms", d)
	}
}

func TestRealClockAfter(t *testing.T) {
	rc := RealClock{}
	select {
	case <-rc.After(10 * time.Millisecond):
	case <-time.After(time.Second):
		t.Fatal("After() channel did not fire in time")
	}
}

func TestRealClockSleep(t *testing.T) {
	rc := RealClock{}
	start := time.Now()
	rc.Sleep(10 * time.Millisecond)
	if elapsed := time.Since(start); elapsed < 10*time.Millisecond {
		t.Fatalf("Sleep() returned after %v, want at least 10ms", elapsed)
	}
}

func TestRealClockTick(t *testing.T) {
	rc := RealClock{}
	ch := rc.Tick(10 * time.Millisecond)
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("Tick() channel did not fire in time")
	}
}

func TestRealClockNewTimer(t *testing.T) {
	rc := RealClock{}
	timer := rc.NewTimer(10 * time.Millisecond)
	select {
	case <-timer.C():
	case <-time.After(time.Second):
		t.Fatal("NewTimer() channel did not fire in time")
	}
	if timer.Reset(time.Minute) {
		t.Fatal("Reset() on an already-fired timer should return false")
	}
	if !timer.Stop() {
		t.Fatal("Stop() right after Reset() should return true")
	}
}

func TestRealClockNewTimerStop(t *testing.T) {
	rc := RealClock{}
	timer := rc.NewTimer(time.Minute)
	if !timer.Stop() {
		t.Fatal("Stop() on a fresh timer should return true")
	}
	select {
	case <-timer.C():
		t.Fatal("stopped timer should not fire")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestRealClockAfterFunc(t *testing.T) {
	rc := RealClock{}
	done := make(chan struct{})
	timer := rc.AfterFunc(10*time.Millisecond, func() { close(done) })
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("AfterFunc() callback did not run in time")
	}
	if timer.Stop() {
		t.Fatal("Stop() after AfterFunc already ran should return false")
	}
}

func TestRealClockNewTicker(t *testing.T) {
	rc := RealClock{}
	ticker := rc.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	select {
	case <-ticker.C():
	case <-time.After(time.Second):
		t.Fatal("NewTicker() channel did not fire in time")
	}
}
