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
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type MakefileMigrator struct{}

func (m MakefileMigrator) Migrate(workingDir string) error {
	// Find Makefile in current or parent directories
	makefilePath, makefile, err := readFile(workingDir, "Makefile")
	if err != nil {
		return fmt.Errorf("could not find Makefile: %v", err)
	}

	connectorName := filepath.Base(workingDir)
	old := fmt.Sprintf(`build:
	go build -ldflags "-X 'github.com/conduitio/%s.version=${VERSION}'" -o %s cmd/connector/main.go`, connectorName, connectorName)

	newBuild := fmt.Sprintf(`build:
	sed -i '/specification:/,/version:/ s/version: .*/version: '"${VERSION}"'/' connector.yaml
	go build -o %s cmd/connector/main.go`, connectorName)
	// Replace the build target
	makefile = strings.ReplaceAll(makefile, old, newBuild)

	// Write back to file
	err = os.WriteFile(makefilePath, []byte(makefile), 0644)
	if err != nil {
		return fmt.Errorf("could not write to %s: %v", makefilePath, err)
	}

	fmt.Printf("Updated build target in %s\n", makefilePath)
	return nil
}
