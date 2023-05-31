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

package keymutex

import (
	"context"
)

// KeyMutex is a thread-safe interface for acquiring locks on arbitrary strings.
type KeyMutex interface {
	// Acquires a lock associated with the specified ID, creates the lock if one doesn't already exist.
	LockKey(id string)

	// Acquires a lock associated with the specified ID, creates the lock if one doesn't already exist.
	// This method will return true if the lock was acquired successfully. It will return false
	// if the context was cancelled before the lock could be acquired.
	LockKeyWithContext(id string, ctx context.Context) bool

	// Releases the lock associated with the specified ID.
	// Returns an error if the specified ID doesn't exist.
	UnlockKey(id string) error
}
