// Copyright © 2024 Meroxa, Inc.
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

package internal

import (
	"os"
	"path/filepath"
)

type Migrator interface {
	Migrate(workingDir string) error
}

// readFile reads filePath in workingDir. Returns: path, contents, error
func readFile(workingDir, filePath string) (string, string, error) {
	p := filepath.Join(workingDir, filePath)

	// Check if file exists
	_, err := os.Stat(p)
	if err != nil {
		return "", "", err
	}

	contents, err := os.ReadFile(p)
	if err != nil {
		return "", "", err
	}

	return p, string(contents), nil
}
