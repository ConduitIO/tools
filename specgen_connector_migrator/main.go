package main

import (
	"fmt"
	"log"
	"os"

	"github.com/conduitio-labs/connector-migrator/internal"
)

func main() {
	// Working directory can be passed as an argument or use current directory
	workingDir := "."
	if len(os.Args) > 1 {
		workingDir = os.Args[1]
	}

	migrators := []internal.Migrator{
		internal.CheckoutNewBranch{},
		internal.ToolsGo{},
		internal.UpgradeSDK{},
		internal.ConnectorGoMigrator{},
		internal.UpdateSourceGo{},
		internal.UpdateDestinationGo{},
		internal.WriteConnectorYaml{},
		internal.MakefileMigrator{},
		internal.DeletedParamGen{},
	}

	for _, m := range migrators {
		fmt.Printf("Running %T\n\n", m)
		err := m.Migrate(workingDir)
		if err != nil {
			log.Fatalf("%T failed: %v", m, err)
		}
		fmt.Printf("\nDone with %T\n-----------\n", m)
	}
}
