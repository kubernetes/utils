/*
Copyright 2017 The Kubernetes Authors.

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

package io

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// ErrLimitReached means that the read limit is reached.
var ErrLimitReached = errors.New("the read limit is reached")

// ConsistentRead repeatedly reads a file until it gets the same content twice.
// This is useful when reading files in /proc that are larger than page size
// and kernel may modify them between individual read() syscalls.
// It returns InconsistentReadError when it cannot get a consistent read in
// given nr. of attempts. Caller should retry, kernel is probably under heavy
// mount/unmount load.
func ConsistentRead(filename string, attempts int) ([]byte, error) {
	return consistentReadSync(filename, attempts, nil)
}

// consistentReadSync is the main functionality of ConsistentRead but
// introduces a sync callback that can be used by the tests to mutate the file
// from which the test data is being read
func consistentReadSync(filename string, attempts int, sync func(int)) ([]byte, error) {

	// get the size of the file.
	size, err := getFileSize(filename)
	if err != nil {
		return nil, err
	}

	// if the size of the file is less than 512 bytes, assume that reading from
	// /proc or /sys like cases for Linux.
	if size < 512 {
		if contents, err := consistentReadFileNoStat(filename); err == nil {
			return contents, nil
		}
		// fall back to consistent read because this is now a utility and an error here can break
		// expected behavior.
	}

	oldContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	for i := 0; i < attempts; i++ {
		if sync != nil {
			sync(i)
		}
		newContent, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		if bytes.Compare(oldContent, newContent) == 0 {
			return newContent, nil
		}
		// Files are different, continue reading
		oldContent = newContent
	}
	return nil, InconsistentReadError{filename, attempts}
}

// get filesize returns the filesize
func getFileSize(filename string) (int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}

	var size int
	// Calling stat should be okay as it supports both Windows and Linux
	// Ref (unix and posix compliant): https://cs.opensource.google/go/go/+/refs/tags/go1.17.1:src/os/stat_unix.go;l=16
	// Ref (windows): https://cs.opensource.google/go/go/+/refs/tags/go1.17.1:src/os/stat_windows.go
	if info, err := f.Stat(); err == nil {
		size64 := info.Size()
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	}
	size++ // one byte for final read at EOF
	return size, nil
}

// consistentReadFileNoStat is inspired from prometheus code base for /proc and /sys
func consistentReadFileNoStat(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()


	const maxBufferSize = 1024 * 512
	var contents = make([]byte, maxBufferSize)

	// use f.Read as it points to File.Read which will directly call os dependent functions
	// For Linux: https://cs.opensource.google/go/go/+/master:src/syscall/syscall_unix.go;l=188;bpv=1;bpt=1?q=syscall.Read&ss=go%2Fgo
	// Also, do not prefer ReadAtMost as it internally results in a loop for read depending on the buffer capacity.
	// This will result in multiple syscall.Read
	n, err := f.Read(contents)
	if err != nil {
		// not returning contents to ensure that the bytes are not used up the stack as they occupy memory.
		return []byte{}, err
	}
	// slice the byte before returning to the x number of bytes so that the maximum buffer is not reached
	return contents[:n], nil
}

// InconsistentReadError is returned from ConsistentRead when it cannot get
// a consistent read in given nr. of attempts. Caller should retry, kernel is
// probably under heavy mount/unmount load.
type InconsistentReadError struct {
	filename string
	attempts int
}

func (i InconsistentReadError) Error() string {
	return fmt.Sprintf("could not get consistent content of %s after %d attempts", i.filename, i.attempts)
}

var _ error = InconsistentReadError{}

func IsInconsistentReadError(err error) bool {
	if _, ok := err.(InconsistentReadError); ok {
		return true
	}
	return false
}

// ReadAtMost reads up to `limit` bytes from `r`, and reports an error
// when `limit` bytes are read.
func ReadAtMost(r io.Reader, limit int64) ([]byte, error) {
	limitedReader := &io.LimitedReader{R: r, N: limit}
	data, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		return data, err
	}
	if limitedReader.N <= 0 {
		return data, ErrLimitReached
	}
	return data, nil
}
