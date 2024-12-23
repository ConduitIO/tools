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
	"os/exec"
)

type CheckoutNewBranch struct {
}

func (c CheckoutNewBranch) Migrate(workingDir string) error {
	// Create command executor in the specified directory
	cmd := exec.Command("git", "checkout", "main")
	cmd.Dir = workingDir

	// Checkout main
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to checkout main: %v, output: %s", err, out)
	}

	// Pull latest
	cmd = exec.Command("git", "pull", "origin", "main")
	cmd.Dir = workingDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to pull: %v, output: %s", err, out)
	}

	// Create new branch
	cmd = exec.Command("git", "checkout", "-b", "haris/specgen")
	cmd.Dir = workingDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create branch: %v, output: %s", err, out)
	}

	// Execute 'git add .'
	addCmd := exec.Command("git", "add", ".")
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed executing git add: %v\n", err)
	}

	// Execute 'git commit -am "Generate connector.yaml"'
	commitCmd := exec.Command("git", "commit", "-am", "Generate connector.yaml")
	if output, err := commitCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed executing git commit: %v\nOutput: %s\n", err, output)
	}

	fmt.Println("Successfully added and committed changes")
	return nil
}