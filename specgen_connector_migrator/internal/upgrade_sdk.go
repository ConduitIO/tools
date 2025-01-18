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
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type UpgradeSDK struct {
}

func (u UpgradeSDK) Migrate(workingDir string) error {
	module := "github.com/conduitio/conduit-connector-sdk"
	version := "main"

	// Run go get command
	err := runCommand(workingDir, "go", "get", fmt.Sprintf("%s@%s", module, version))
	if err != nil {
		return fmt.Errorf("could not run `go get`: %w", err)
	}

	// Run go mod tidy
	err = runCommand(workingDir, "go", "mod", "tidy")
	if err != nil {
		return fmt.Errorf("Failed to run go mod tidy: %v\n", err)
	}

	return nil
}

func runCommand(workingDir string, command string, args ...string) error {
	// Construct the command
	cmd := exec.Command(command, args...)

	// Set the working directory
	cmd.Dir = workingDir

	// Capture stdout and stderr
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// Set environment variables from the current process
	cmd.Env = os.Environ()

	// Run the command
	err := cmd.Run()

	// Print outputs
	if outBuf.Len() > 0 {
		fmt.Println(command + " STDOUT:")
		fmt.Println(outBuf.String())
	}

	if errBuf.Len() > 0 {
		fmt.Println(command + " STDERR:")
		fmt.Println(errBuf.String())
	}

	// Return any error that occurred during command execution
	if err != nil {
		return fmt.Errorf("error running %s: %v", command, err)
	}

	return nil
}
