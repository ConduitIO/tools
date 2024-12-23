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
	"regexp"
	"strings"
)

type ConnectorGoMigrator struct {
}

func (a ConnectorGoMigrator) Migrate(workingDir string) error {
	connectorGoPath, connectorGo, err := readFile(workingDir, "connector.go")
	if err != nil {
		return err
	}

	updatedConnectorGo := strings.ReplaceAll(
		connectorGo,
		`// limitations under the License.`,
		`// limitations under the License.

//go:generate specgen`,
	)

	// import embed
	updatedConnectorGo = strings.ReplaceAll(connectorGo,
		`import (
	sdk "github.com/conduitio/conduit-connector-sdk"
)`,
		`import (
	_ "embed"

	sdk "github.com/conduitio/conduit-connector-sdk"
)

//go:embed connector.yaml
var specs string

`)
	// import embed
	updatedConnectorGo = strings.ReplaceAll(connectorGo,
		`import sdk "github.com/conduitio/conduit-connector-sdk"`,
		`import (
	_ "embed"

	sdk "github.com/conduitio/conduit-connector-sdk"
)

//go:embed connector.yaml
var specs string

`)

	// Compile the regex pattern
	regex := regexp.MustCompile(`NewSpecification:.*`)

	// Replace the line with the new specification
	updatedConnectorGo = regex.ReplaceAllString(updatedConnectorGo, "NewSpecification: sdk.YAMLSpecification(specs),")

	// Write back to file
	err = os.WriteFile(connectorGoPath, []byte(updatedConnectorGo), 0644)
	if err != nil {
		return fmt.Errorf("could not write to %s: %v", connectorGoPath, err)
	}

	fmt.Printf("Updated connector.go target in %s\n", connectorGoPath)

	return nil
}
