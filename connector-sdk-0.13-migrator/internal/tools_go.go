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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ToolsGo struct {
}

func (t ToolsGo) Migrate(workingDir string) error {
	toolsGoPath, toolsGo, err := readFile(workingDir, "tools.go")
	if errors.Is(err, os.ErrNotExist) {
		return t.migrateToolsDir(workingDir)
	}

	if err != nil {
		return err
	}

	updatedToolsGo := strings.ReplaceAll(
		toolsGo,
		"_ \"github.com/conduitio/conduit-commons/paramgen\"",
		"_ \"github.com/conduitio/conduit-connector-sdk/conn-sdk-cli\"",
	)

	err = os.WriteFile(toolsGoPath, []byte(updatedToolsGo), 0644)
	if err != nil {
		return fmt.Errorf("failed writing new contents of tools.go: %w", err)
	}

	return nil
}

func (t ToolsGo) migrateToolsDir(workingDir string) error {
	toolsGoPath, toolsGo, err := readFile(workingDir, "tools/go.mod")
	if err != nil {
		return fmt.Errorf("failed reading tools/go.mod: %w", err)
	}
	updatedGoMod := strings.ReplaceAll(
		toolsGo,
		"github.com/conduitio/conduit-commons/paramgen",
		"github.com/conduitio/conduit-connector-sdk/conn-sdk-cli",
	)

	err = os.WriteFile(toolsGoPath, []byte(updatedGoMod), 0644)
	if err != nil {
		return fmt.Errorf("failed writing new contents of tools/go.mod: %w", err)
	}

	err = runGoModTidy(filepath.Join(workingDir, "tools"))
	if err != nil {
		return fmt.Errorf("failed to run go mod tidy in tools directory: %w", err)
	}

	return nil
}

func runGoModTidy(dir string) error {
	// Convert to absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if directory exists
	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", absDir)
	}

	// Check if go.mod exists in the directory
	goModPath := filepath.Join(absDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found in directory: %s", absDir)
	}

	// Create the command
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = absDir

	// Set up output capture
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("Successfully ran 'go mod tidy' in %s\n", absDir)
	if len(output) > 0 {
		fmt.Printf("Output: %s\n", string(output))
	}

	return nil
}
