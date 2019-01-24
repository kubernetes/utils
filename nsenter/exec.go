// +build linux

/*
Copyright 2018 The Kubernetes Authors.

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

package nsenter

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/go-logr/logr"
	"k8s.io/utils/exec"
)

// Executor wraps executor interface to be executed via nsenter
type Executor struct {
	// Exec implementation
	executor exec.Interface
	// Path to the host's root proc path
	hostProcMountNsPath string
	// How to log
	log logr.InfoLogger
}

// NewNsenterExecutor returns new nsenter based executor.  If logging is enabled
// (see WithLogger) Executor will log executions, but will return errors without
// logging them.
func NewNsenterExecutor(hostRootFsPath string, executor exec.Interface) *Executor {
	hostProcMountNsPath := filepath.Join(hostRootFsPath, mountNsPath)
	nsExecutor := &Executor{
		hostProcMountNsPath: hostProcMountNsPath,
		executor:            executor,
		log:                 nil,
	}
	return nsExecutor
}

// WithLogger returns the same executor, but configures it for logging.
func (e *Executor) WithLogger(log logr.InfoLogger) *Executor {
	e.log = log
	return e
}

// Command returns a command wrapped with nsenter.
func (e *Executor) Command(cmd string, args ...string) exec.Cmd {
	fullArgs := append([]string{fmt.Sprintf("--mount=%s", e.hostProcMountNsPath), "--"},
		append([]string{cmd}, args...)...)
	if e.log != nil {
		e.log.Info("Running nsenter", "bin", nsenterPath, "args", fullArgs)
	}
	return e.executor.Command(nsenterPath, fullArgs...)
}

// CommandContext returns a CommandContext wrapped with nsenter.
func (e *Executor) CommandContext(ctx context.Context, cmd string, args ...string) exec.Cmd {
	fullArgs := append([]string{fmt.Sprintf("--mount=%s", e.hostProcMountNsPath), "--"},
		append([]string{cmd}, args...)...)
	if e.log != nil {
		e.log.Info("Running nsenter", "bin", nsenterPath, "args", fullArgs)
	}
	return e.executor.CommandContext(ctx, nsenterPath, fullArgs...)
}

// LookPath returns a LookPath wrapped with nsenter.
func (*Executor) LookPath(file string) (string, error) {
	return "", fmt.Errorf("LookPath() is not supported for nsenter.Executor: %s", file)
}
