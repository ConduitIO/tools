package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/conduitio-labs/connector-migrator/internal"
)

var allMigrators = []internal.Migrator{
	internal.ToolsGo{},
	internal.UpgradeSDK{},
	internal.ConnectorGoMigrator{},
	internal.UpdateSourceGo{},
	internal.UpdateDestinationGo{},
	internal.WriteConnectorYaml{},
	internal.DeleteParamGen{},
	internal.DeleteSpecGo{},
	internal.WorkflowRelease{},
	internal.GoReleaserMigrator{},
	internal.MakefileMigrator{},
	internal.ScriptsMigrator{},
}

func main() {
	// Working directory can be passed as an argument or use current directory
	workingDir := "."
	if len(os.Args) > 1 {
		workingDir = os.Args[1]
	}
	migrator := ""
	if len(os.Args) > 2 {
		migrator = os.Args[2]
	}

	var migrators []internal.Migrator
	if migrator == "" {
		migrators = allMigrators
	} else {
		for _, m := range allMigrators {
			if reflect.TypeOf(m).Name() == migrator {
				migrators = []internal.Migrator{m}
				break
			}
		}
	}

	fmt.Printf("Migrating %v\n", workingDir)

	for _, m := range migrators {
		fmt.Printf("Running %T\n\n", m)
		err := m.Migrate(workingDir)
		if err != nil {
			log.Fatalf("%T failed: %v", m, err)
		}
		fmt.Printf("\nDone with %T\n-----------\n", m)
	}
}
