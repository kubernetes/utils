/*
Copyright 2015 The Kubernetes Authors.

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

package trace

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"k8s.io/klog/v2"
)

// Field is a key value pair that provides additional details about the trace.
type Field struct {
	Key   string
	Value interface{}
}

type stepTrace interface {
	time() time.Time
	writeStep(b *bytes.Buffer, formatter string, stepDuration time.Duration, lastStepTime time.Time) time.Time
}

func (t *Trace) sortStepTrace() []stepTrace {
	var arr []stepTrace

	for _, step := range t.steps {
		arr = append(arr, step)
	}

	for _, trace := range t.nestedTrace {
		arr = append(arr, trace)
	}

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].time().Before(arr[j].time())
	})

	return arr
}

func (f Field) format() string {
	return fmt.Sprintf("%s:%v", f.Key, f.Value)
}

func writeFields(b *bytes.Buffer, l []Field) {
	for i, f := range l {
		b.WriteString(f.format())
		if i < len(l)-1 {
			b.WriteString(",")
		}
	}
}

func writeMainLog(b *bytes.Buffer, msg string, totalTime time.Duration, startTime time.Time, fields []Field) {
	b.WriteString(fmt.Sprintf("%q ", msg))
	if len(fields) > 0 {
		writeFields(b, fields)
		b.WriteString(" ")
	}

	b.WriteString(fmt.Sprintf("%vms (%v)", durationToMilliseconds(totalTime), startTime.Format("15:04:00.000")))
}

func durationToMilliseconds(timeDuration time.Duration) int64 {
	return timeDuration.Nanoseconds() / 1e6
}

type traceStep struct {
	stepTime time.Time
	msg      string
	fields   []Field
}

func (s traceStep) time() time.Time {
	return s.stepTime
}

func (s traceStep) writeStep(b *bytes.Buffer, formatter string, stepThreshold time.Duration,
	lastStepTime time.Time) time.Time {
	stepDuration := s.stepTime.Sub(lastStepTime)
	if stepThreshold == 0 || stepDuration > stepThreshold || klog.V(4).Enabled() {
		b.WriteString(fmt.Sprintf("%v---", formatter))
		writeMainLog(b, s.msg, stepDuration, s.stepTime, s.fields)
	}
	return s.stepTime
}

// Trace keeps track of a set of "steps" and allows us to log a specific
// step if it took longer than its share of the total allowed time
type Trace struct {
	name        string
	fields      []Field
	startTime   time.Time
	steps       []traceStep
	nestedTrace []*Trace
}

func (t *Trace) time() time.Time {
	return t.startTime
}

func (t *Trace) writeStep(b *bytes.Buffer, formatter string, stepThreshold time.Duration,
	lastStepTime time.Time) time.Time {
	b.WriteString(fmt.Sprintf("%v[", formatter))
	writeMainLog(b, t.name, t.TotalTime(), t.startTime, t.fields)
	_ = writeTrace(b, t, formatter+" ", stepThreshold)
	b.WriteString("]")
	return time.Time{}
}

// New creates a Trace with the specified name. The name identifies the operation to be traced. The
// Fields add key value pairs to provide additional details about the trace, such as operation inputs.
func New(name string, fields ...Field) *Trace {
	return &Trace{name: name, startTime: time.Now(), fields: fields}
}

// Step adds a new step with a specific message. Call this at the end of an execution step to record
// how long it took. The Fields add key value pairs to provide additional details about the trace
// step.
func (t *Trace) Step(msg string, fields ...Field) {
	if t.steps == nil {
		// traces almost always have less than 6 steps, do this to avoid more than a single allocation
		t.steps = make([]traceStep, 0, 6)
	}
	t.steps = append(t.steps, traceStep{stepTime: time.Now(), msg: msg, fields: fields})
}

// Nest adds a nested trace with the given message and fields and returns it.
func (t *Trace) Nest(msg string, fields ...Field) *Trace {
	newTrace := New(msg, fields...)
	t.nestedTrace = append(t.nestedTrace, newTrace)
	return newTrace
}

// Log is used to dump all the steps in the Trace. It also logs the nested trace messages using indentation
func (t *Trace) Log() {
	// an explicit logging request should dump all the steps out at the higher level
	t.logWithStepThreshold(0)
}

func (t *Trace) logWithStepThreshold(stepThreshold time.Duration) {
	var buffer bytes.Buffer
	tracenum := rand.Int31()
	endTime := time.Now()

	totalTime := endTime.Sub(t.startTime)
	buffer.WriteString(fmt.Sprintf("Trace[%d]: %q ", tracenum, t.name))
	if len(t.fields) > 0 {
		writeFields(&buffer, t.fields)
		buffer.WriteString(" ")
	}
	buffer.WriteString(fmt.Sprintf("(%v) (total time: %vms):", t.startTime.Format("02-Jan-2006 15:04:00.000"), totalTime.Milliseconds()))
	lastStepTime := writeTrace(&buffer, t, fmt.Sprintf("\nTrace[%d]: ", tracenum), stepThreshold)
	stepDuration := endTime.Sub(lastStepTime)
	if stepThreshold == 0 || stepDuration > stepThreshold || klog.V(4).Enabled()  {
		buffer.WriteString(fmt.Sprintf("\nTrace[%d]: [%v] [%v] END\n", tracenum, endTime.Sub(t.startTime), stepDuration))
	}

	klog.Info(buffer.String())
}

func writeTrace(b *bytes.Buffer, t *Trace, formatter string, stepThreshold time.Duration) time.Time {
	lastStepTime := t.startTime
	stepAndTraces := t.sortStepTrace()
	if len(stepAndTraces) == 0 {
		return lastStepTime
	}
	for _, stepOrTrace := range stepAndTraces {
		blankTime := time.Time{}
		stepTime := stepOrTrace.writeStep(b, formatter, stepThreshold, lastStepTime)
		if stepTime != blankTime {
			lastStepTime = stepTime
		}
	}
	return lastStepTime
}

// LogIfLong is used to dump steps that took longer than its share
func (t *Trace) LogIfLong(threshold time.Duration) {
	if time.Since(t.startTime) >= threshold {
		// if any step took more than it's share of the total allowed time, it deserves a higher log level
		stepThreshold := threshold / time.Duration(len(t.steps)+1)
		t.logWithStepThreshold(stepThreshold)
	}
}

// TotalTime can be used to figure out how long it took since the Trace was created
func (t *Trace) TotalTime() time.Duration {
	return time.Since(t.startTime)
}
