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
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed release.yaml
var workflowFiles embed.FS

type WorkflowRelease struct{}

func (w WorkflowRelease) Migrate(workingDir string) error {
	// Check for both possible file extensions
	possiblePaths := []string{
		filepath.Join(workingDir, ".github", "workflows", "release.yml"),
		filepath.Join(workingDir, ".github", "workflows", "release.yaml"),
	}

	// Check if either file exists
	var existingPath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			existingPath = path
			break
		}
	}

	// Read the embedded workflow file
	workflowContent, err := workflowFiles.ReadFile("release.yaml")
	if err != nil {
		return fmt.Errorf("failed to read workflow file: %w", err)
	}

	// Create the target directory if it doesn't exist
	targetDir := filepath.Join(workingDir, ".github", "workflows")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create workflows directory: %w", err)
	}

	// Write the new file
	if err := os.WriteFile(existingPath, workflowContent, 0644); err != nil {
		return fmt.Errorf("failed to write workflow file: %w", err)
	}

	return nil
}
