// Copyright © 2025 Meroxa, Inc.
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
	// Construct path to Makefile
	makefilePath := filepath.Join(workingDir, "Makefile")

	// Read the Makefile content
	content, err := os.ReadFile(makefilePath)
	if err != nil {
		return fmt.Errorf("failed to read Makefile: %w", err)
	}

	// Define the old and new content with proper tab indentation
	oldContent := `.PHONY: generate
generate:
   go generate ./...`

	newContent := `.PHONY: generate
generate:
   go generate ./...
   conn-sdk-cli readmegen -w`

	// Perform the replacement
	updatedContent := strings.ReplaceAll(string(content), oldContent, newContent)

	// Write the modified content back to the file
	err = os.WriteFile(makefilePath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated Makefile: %w", err)
	}

	return nil
}
