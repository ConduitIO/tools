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
	"os"
	"strings"
)

type ToolsGo struct {
}

func (t ToolsGo) Migrate(workingDir string) error {
	toolsGoPath, toolsGo, err := readFile(workingDir, "tools.go")
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
		return err
	}

	return nil
}
