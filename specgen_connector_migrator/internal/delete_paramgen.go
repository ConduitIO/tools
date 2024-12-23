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
	"regexp"
	"strings"
)

type DeletedParamGen struct {
}

func (d DeletedParamGen) Migrate(workingDir string) error {
	// Walk through the directory recursively
	err := filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		// Check if there was an error accessing the path
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Read entire file into string
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if strings.Contains(string(content), "// Code generated by paramgen. DO NOT EDIT.") {
			fmt.Printf("deleting %s\n", path)

			err := os.Remove(file.Name())
			if err != nil {
				return fmt.Errorf("removing file %s: %w", file, err)
			}
			return nil
		}

		if strings.Contains(string(content), "//go:generate paramgen") {
			fmt.Printf("updating %s (removing paramgen)\n", path)
			// Find and replace the generate line (and its newline)
			regex := regexp.MustCompile(`(?m)^[\t ]*//go:generate paramgen.*\n`)
			newContent := regex.ReplaceAllString(string(content), "")

			// Write back to file
			return os.WriteFile(path, []byte(newContent), 0644)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}