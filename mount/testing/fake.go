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

package testingmount

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/golang/glog"

	"k8s.io/utils/mount"
)

// Fake for testing.
type FakeExec struct {
	runHook runHook
}
type runHook func(cmd string, args ...string) ([]byte, error)

func (f *FakeExec) Run(cmd string, args ...string) ([]byte, error) {
	if f.runHook != nil {
		return f.runHook(cmd, args...)
	}
	return nil, nil
}

func NewFakeExec(run runHook) *FakeExec {
	return &FakeExec{runHook: run}
}

// FakeMounter implements mount.Interface for tests.
type FakeMounter struct {
	MountPoints []mount.MountPoint
	Log         []FakeAction
	// Some tests run things in parallel, make sure the mounter does not produce
	// any golang's DATA RACE warnings.
	mutex sync.Mutex
}

var _ mount.Interface = &FakeMounter{}

// Values for FakeAction.Action
const FakeActionMount = "mount"
const FakeActionUnmount = "unmount"

// FakeAction objects are logged every time a fake mount or unmount is called.
type FakeAction struct {
	Action string // "mount" or "unmount"
	Target string // applies to both mount and unmount actions
	Source string // applies only to "mount" actions
	FSType string // applies only to "mount" actions
}

func (f *FakeMounter) ResetLog() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.Log = []FakeAction{}
}

func (f *FakeMounter) Mount(source string, target string, fstype string, options []string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	opts := []string{}

	for _, option := range options {
		// find 'bind' option
		if option == "bind" {
			// This is a bind-mount. In order to mimic linux behaviour, we must
			// use the original device of the bind-mount as the real source.
			// E.g. when mounted /dev/sda like this:
			//      $ mount /dev/sda /mnt/test
			//      $ mount -o bind /mnt/test /mnt/bound
			// then /proc/mount contains:
			// /dev/sda /mnt/test
			// /dev/sda /mnt/bound
			// (and not /mnt/test /mnt/bound)
			// I.e. we must use /dev/sda as source instead of /mnt/test in the
			// bind mount.
			for _, mnt := range f.MountPoints {
				if source == mnt.Path {
					source = mnt.Device
					break
				}
			}
		}
		// find 'ro' option
		if option == "ro" {
			// reuse MountPoint.Opts field to mark mount as readonly
			opts = append(opts, "ro")
		}
	}

	// If target is a symlink, get its absolute path
	absTarget, err := filepath.EvalSymlinks(target)
	if err != nil {
		absTarget = target
	}

	f.MountPoints = append(f.MountPoints, mount.MountPoint{Device: source, Path: absTarget, Type: fstype, Opts: opts})
	glog.V(5).Infof("Fake mounter: mounted %s to %s", source, absTarget)
	f.Log = append(f.Log, FakeAction{Action: FakeActionMount, Target: absTarget, Source: source, FSType: fstype})
	return nil
}

func (f *FakeMounter) Unmount(target string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// If target is a symlink, get its absolute path
	absTarget, err := filepath.EvalSymlinks(target)
	if err != nil {
		absTarget = target
	}

	newMountpoints := []mount.MountPoint{}
	for _, mp := range f.MountPoints {
		if mp.Path == absTarget {
			glog.V(5).Infof("Fake mounter: unmounted %s from %s", mp.Device, absTarget)
			// Don't copy it to newMountpoints
			continue
		}
		newMountpoints = append(newMountpoints, mount.MountPoint{Device: mp.Device, Path: mp.Path, Type: mp.Type})
	}
	f.MountPoints = newMountpoints
	f.Log = append(f.Log, FakeAction{Action: FakeActionUnmount, Target: absTarget})
	return nil
}

func (f *FakeMounter) List() ([]mount.MountPoint, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return f.MountPoints, nil
}

func (f *FakeMounter) IsMountPointMatch(mp mount.MountPoint, dir string) bool {
	return mp.Path == dir
}

func (f *FakeMounter) IsNotMountPoint(dir string) (bool, error) {
	return mount.IsNotMountPoint(f, dir)
}

func (f *FakeMounter) IsLikelyNotMountPoint(file string) (bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	_, err := os.Stat(file)
	if err != nil {
		return true, err
	}

	// If file is a symlink, get its absolute path
	absFile, err := filepath.EvalSymlinks(file)
	if err != nil {
		absFile = file
	}

	for _, mp := range f.MountPoints {
		if mp.Path == absFile {
			glog.V(5).Infof("isLikelyNotMountPoint for %s: mounted %s, false", file, mp.Path)
			return false, nil
		}
	}
	glog.V(5).Infof("isLikelyNotMountPoint for %s: true", file)
	return true, nil
}

func (f *FakeMounter) DeviceOpened(pathname string) (bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	for _, mp := range f.MountPoints {
		if mp.Device == pathname {
			return true, nil
		}
	}
	return false, nil
}

func (f *FakeMounter) PathIsDevice(pathname string) (bool, error) {
	return true, nil
}

func (f *FakeMounter) GetDeviceNameFromMount(mountPath, pluginDir string) (string, error) {
	return getDeviceNameFromMount(f, mountPath, pluginDir)
}

func (f *FakeMounter) MakeRShared(path string) error {
	return nil
}

func (f *FakeMounter) GetFileType(pathname string) (mount.FileType, error) {
	return mount.FileType("fake"), nil
}

func (f *FakeMounter) MakeDir(pathname string) error {
	return nil
}

func (f *FakeMounter) MakeFile(pathname string) error {
	return nil
}

func (f *FakeMounter) ExistsPath(pathname string) (bool, error) {
	return false, errors.New("not implemented")
}

func (f *FakeMounter) PrepareSafeSubpath(subPath mount.Subpath) (newHostPath string, cleanupAction func(), err error) {
	return subPath.Path, nil, nil
}

func (f *FakeMounter) CleanSubPaths(podDir string, volumeName string) error {
	return nil
}
func (mounter *FakeMounter) SafeMakeDir(pathname string, base string, perm os.FileMode) error {
	return nil
}

func (f *FakeMounter) GetMountRefs(pathname string) ([]string, error) {
	realpath, err := filepath.EvalSymlinks(pathname)
	if err != nil {
		// Ignore error in FakeMounter, because we actually didn't create files.
		realpath = pathname
	}
	return getMountRefsByDev(f, realpath)
}

// getMountRefsByDev finds all references to the device provided
// by mountPath; returns a list of paths.
// Note that mountPath should be path after the evaluation of any symblolic links.
func getMountRefsByDev(mounter mount.Interface, mountPath string) ([]string, error) {
	mps, err := mounter.List()
	if err != nil {
		return nil, err
	}

	// Finding the device mounted to mountPath
	diskDev := ""
	for i := range mps {
		if mountPath == mps[i].Path {
			diskDev = mps[i].Device
			break
		}
	}

	// Find all references to the device.
	var refs []string
	for i := range mps {
		if mps[i].Device == diskDev || mps[i].Device == mountPath {
			if mps[i].Path != mountPath {
				refs = append(refs, mps[i].Path)
			}
		}
	}
	return refs, nil
}

func (f *FakeMounter) GetFSGroup(pathname string) (int64, error) {
	return -1, errors.New("GetFSGroup not implemented")
}

func (f *FakeMounter) GetSELinuxSupport(pathname string) (bool, error) {
	return false, errors.New("GetSELinuxSupport not implemented")
}

func (f *FakeMounter) GetMode(pathname string) (os.FileMode, error) {
	return 0, errors.New("not implemented")
}
