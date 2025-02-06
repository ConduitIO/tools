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

//go:embed scripts/*
var scripts embed.FS

type ScriptsMigrator struct{}

func (s ScriptsMigrator) Migrate(workingDir string) error {
	// Create the destination directory if it doesn't exist
	destDir := filepath.Join(workingDir, "scripts")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	srcDir := "scripts"
	entries, err := scripts.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read scripts directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		content, err := scripts.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", srcPath, err)
		}

		destPath := filepath.Join(destDir, entry.Name())
		err = os.WriteFile(destPath, content, 0755)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}
	}

	return nil
}
