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
	"regexp"
)

type GoReleaserMigrator struct{}

func (g GoReleaserMigrator) Migrate(workingDir string) error {
	// Try both yaml extensions
	possibleFiles := []string{".goreleaser.yml", ".goreleaser.yaml"}
	var configPath string

	// Find the first existing config file
	for _, fileName := range possibleFiles {
		path := filepath.Join(workingDir, fileName)
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}
	}

	if configPath == "" {
		return fmt.Errorf("no .goreleaser configuration file found in %s", workingDir)
	}

	// Read the file content
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Convert to string for easier manipulation
	fileContent := string(content)

	// Define the pattern to match and remove
	pattern := `ldflags:\n\s*- "-s -w -X 'github.com/conduitio/conduit-connector-connectorname.version={{ .Tag }}'"`

	// Remove the pattern and any trailing whitespace/newlines
	newContent := regexp.MustCompile(pattern).ReplaceAllString(fileContent, "")

	// Clean up any empty lines that might have been left
	newContent = regexp.MustCompile(`\n\s*\n\s*\n`).ReplaceAllString(newContent, "\n\n")

	// Write the modified content back to the file
	err = os.WriteFile(configPath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated config file: %w", err)
	}

	return nil
}
