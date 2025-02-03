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
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type ScriptsMigrator struct{}

func (s ScriptsMigrator) Migrate(workingDir string) error {
	scriptsDir := "internal/scripts"

	// Check if source directory exists
	if _, err := os.Stat(scriptsDir); os.IsNotExist(err) {
		return fmt.Errorf("scripts directory does not exist: %w", err)
	}

	// Create the destination directory if it doesn't exist
	destDir := filepath.Join(workingDir, "scripts")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Walk through the scripts directory
	return filepath.Walk(scriptsDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %q: %w", path, err)
		}

		// Calculate relative path to maintain directory structure
		relPath, err := filepath.Rel(scriptsDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			// Create directory in destination
			return os.MkdirAll(destPath, 0755)
		}

		// Copy file contents
		return copyFile(path, destPath)
	})
}

// Helper function to copy individual files
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy the contents
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Copy file mode
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	return os.Chmod(dst, sourceInfo.Mode())
}
