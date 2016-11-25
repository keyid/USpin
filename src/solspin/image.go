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

package main

// StartImageBuild will perform all steps up until the point where it is time
// for the pkg.Manager to step in and populate the rootfs.
func (s *SolSpin) StartImageBuild() error {
	var err error

	s.logImage.Info("Preparing workspace")
	if err = s.builder.PrepareWorkspace(); err != nil {
		s.logImage.Error(err)
		return err
	}

	s.logImage.Info("Creating storage")
	if err = s.builder.CreateStorage(); err != nil {
		s.logImage.Error(err)
		return err
	}

	s.logImage.Info("Mounting storage")
	if err = s.builder.MountStorage(); err != nil {
		s.logImage.Error(err)
		return err
	}
	return nil
}
