//
// Copyright © 2016 Ikey Doherty <ikey@solus-project.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package libimage

import (
	"strings"
	"syscall"
	"time"
)

const (
	// UmountMaxTries is the maximum number of times to try unmounting before
	// resorting to lazy detaches
	UmountMaxTries = 3

	// UmountRetryTime is the length of time to wait in between umounts
	UmountRetryTime = 500 * time.Millisecond
)

// A MountEntry is tracked by the MountManager to enable proper cleanup takes
// place
type MountEntry struct {
	SourcePath string // The source of the mount
	MountPoint string // The destination mount point
}

// Umount will attempt to unmount the given path
func (m *MountEntry) Umount() error {
	return syscall.Unmount(m.MountPoint, 0)
}

// UmountForce will attempt to forcibly detach the mountpoint
func (m *MountEntry) UmountForce() error {
	return syscall.Unmount(m.MountPoint, syscall.MNT_FORCE)
}

// UmountLazy will attempt a lazy detach of the node
func (m *MountEntry) UmountLazy() error {
	return syscall.Unmount(m.MountPoint, syscall.MNT_DETACH)
}

// UmountSync will attempt everything possible to umount itself
func (m *MountEntry) UmountSync() error {
	for i := 0; i < UmountMaxTries; i++ {
		if err := m.Umount(); err == nil {
			return nil
		}
		time.Sleep(UmountRetryTime)
	}
	// Still didn't manage to umount it
	if err := m.UmountForce(); err == nil {
		return nil
	}
	return m.UmountLazy()
}

// A MountManager is used to mount and unmount filesystems, and to track them
// so that they are all properly torn down
type MountManager struct {
	mounts map[string]*MountEntry
}

var mountManager *MountManager

func init() {
	mountManager = &MountManager{}
	mountManager.mounts = make(map[string]*MountEntry)
}

// GetMountManager will return the global mount manager
func GetMountManager() *MountManager {
	return mountManager
}

// MountPath will attempt to mount the given sourcepath at the destpath
func (m *MountManager) MountPath(sourcepath, destpath, filesystem string, flags uintptr, options ...string) error {
	optString := ""
	if len(options) > 1 {
		optString = strings.Join(options, ",")
	}
	er := syscall.Mount(sourcepath, destpath, filesystem, flags, optString)
	return er
}

// BindMount will attempt to mount the given sourcepath at the destpath with a binding
func (m *MountManager) BindMount(sourcepath, destpath, filesystem string, options ...string) error {
	return m.MountPath(sourcepath, destpath, filesystem, syscall.MS_BIND, options...)
}

// Mount will attempt to mount the given sourcepath at the destpath with default options
func (m *MountManager) Mount(sourcepath, destpath, filesystem string, options ...string) error {
	return m.MountPath(sourcepath, destpath, filesystem, 0, options...)
}
