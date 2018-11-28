// Copyright 2018 Foo Coders (www.foocoders.io).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"os"
	"path/filepath"
)

func GetAppPath() string {
	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		return dir
	}

	return ""
}

// Exists returns true/false depending if the file indicated in path, exists or not
func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
