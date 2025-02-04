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
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/conduitio/yaml/v3"
)

type SpecificationInfo struct {
	Version     string
	Name        string
	Summary     string
	Description string
	Author      string
}

type YAMLSpecification struct {
	Version       string `yaml:"version"`
	Specification struct {
		Name        string `yaml:"name"`
		Summary     string `yaml:"summary"`
		Description string `yaml:"description"`
		Version     string `yaml:"version"`
		Author      string `yaml:"author"`
	} `yaml:"specification"`
}

type WriteConnectorYaml struct {
}

func (w WriteConnectorYaml) Migrate(workingDir string) error {
	// Extract specification fields
	spec, err := w.extractSpecificationFields(filepath.Join(workingDir, "spec.go"))
	if err != nil {
		return fmt.Errorf("extract specification fields: %w", err)
	}

	// Convert to YAML structure
	yamlSpec, err := w.convertToYAML(spec)
	if err != nil {
		log.Fatalf("Error converting to YAML: %v", err)
	}

	// Marshal to YAML
	yamlData, err := yaml.Marshal(yamlSpec)
	if err != nil {
		log.Fatalf("Error marshaling YAML: %v", err)
	}

	// Write to file
	err = os.WriteFile(filepath.Join(workingDir, "connector.yaml"), yamlData, 0644)
	if err != nil {
		log.Fatalf("Error writing YAML file: %v", err)
	}

	return nil
}

func (w WriteConnectorYaml) convertToYAML(spec *SpecificationInfo) (*YAMLSpecification, error) {
	yamlSpec := &YAMLSpecification{
		Version: "1.0",
	}

	yamlSpec.Specification.Name = spec.Name
	yamlSpec.Specification.Summary = spec.Summary
	yamlSpec.Specification.Description = spec.Description
	yamlSpec.Specification.Version = spec.Version
	yamlSpec.Specification.Author = spec.Author

	return yamlSpec, nil
}

func (w WriteConnectorYaml) extractSpecificationFields(filename string) (*SpecificationInfo, error) {
	// Create a new file set
	fset := token.NewFileSet()

	// Parse the source file
	file, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("error parsing file: %v", err)
	}

	// Look for the Specification function
	for _, decl := range file.Decls {
		// Check if it's a function declaration
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != "Specification" {
			continue
		}

		// Check if the function returns a struct literal
		returnStmt, ok := funcDecl.Body.List[0].(*ast.ReturnStmt)
		if !ok {
			continue
		}

		// Check if it's a composite literal
		compLit, ok := returnStmt.Results[0].(*ast.CompositeLit)
		if !ok {
			continue
		}

		// Create a struct to store extracted information
		spec := &SpecificationInfo{}

		// Iterate through the struct fields
		for _, elt := range compLit.Elts {
			// Type assert to key-value expression
			kvExpr, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			// Extract key and value
			key, ok := kvExpr.Key.(*ast.Ident)
			if !ok {
				continue
			}

			// Extract string value
			strLit, ok := kvExpr.Value.(*ast.BasicLit)
			if !ok || strLit.Kind != token.STRING {
				continue
			}

			// Remove quotes from the string value
			value := strLit.Value[1 : len(strLit.Value)-1]

			// Populate the struct based on the key
			switch key.Name {
			case "Version":
				spec.Version = value
			case "Name":
				spec.Name = value
			case "Summary":
				spec.Summary = value
			case "Description":
				spec.Description = value
			case "Author":
				spec.Author = value
			}
		}

		return spec, nil
	}

	return nil, fmt.Errorf("no Specification function found")
}
