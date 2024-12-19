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
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type UpdateSourceGo struct{}

func (u UpdateSourceGo) Migrate(workingDir string) error {
	// Find Go source files in the working directory
	files, err := filepath.Glob(filepath.Join(workingDir, "*source.go"))
	if err != nil {
		return fmt.Errorf("failed to find Go files: %w", err)
	}

	for _, filename := range files {
		updated, err := u.maybeUpdateSource(filename)
		if err != nil {
			return fmt.Errorf("failed to update file %v: %w", filename, err)
		}

		if updated {
			fmt.Printf("updated source file %s\n", filename)
			break
		}
	}

	return nil
}

// ImplementsSource checks if a given type specification implements the sdk.Source interface
func (u UpdateSourceGo) implementsSource(file *ast.File, n ast.Node) bool {
	typeSpec, ok := n.(*ast.TypeSpec)
	if !ok {
		return false
	}

	// Ensure it's a struct type
	_, isStruct := typeSpec.Type.(*ast.StructType)
	if !isStruct {
		return false
	}

	// Required methods for sdk.Source interface
	requiredMethods := map[string]bool{
		"Open":     false,
		"Read":     false,
		"Ack":      false,
		"Teardown": false,
	}

	// Inspect methods of the type
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Check if the method belongs to this type
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			continue
		}

		// Get the receiver type name
		var receiverTypeName string
		switch recv := funcDecl.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			if ident, ok := recv.X.(*ast.Ident); ok {
				receiverTypeName = ident.Name
			}
		case *ast.Ident:
			receiverTypeName = recv.Name
		}

		// Check if the method belongs to our type and is a required method
		if receiverTypeName == typeSpec.Name.Name {
			if _, exists := requiredMethods[funcDecl.Name.Name]; exists {
				requiredMethods[funcDecl.Name.Name] = true
			}
		}
	}

	// Check if all required methods are implemented
	for _, implemented := range requiredMethods {
		if !implemented {
			return false
		}
	}

	return true
}

func (u UpdateSourceGo) getMethod(file *ast.File, structNode ast.Node, methodName string) (*ast.FuncDecl, error) {
	// Get the struct name
	typeSpec, ok := structNode.(*ast.TypeSpec)
	if !ok {
		return nil, fmt.Errorf("node is not a TypeSpec")
	}
	structName := typeSpec.Name.Name

	// Look for the method in all declarations
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Check if it's a method (has a receiver) and matches the name
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 || funcDecl.Name.Name != methodName {
			continue
		}

		// Check if the receiver type matches our struct
		var receiverName string
		switch recv := funcDecl.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			// Pointer receiver (*StructName)
			if ident, ok := recv.X.(*ast.Ident); ok {
				receiverName = ident.Name
			}
		case *ast.Ident:
			// Value receiver (StructName)
			receiverName = recv.Name
		}

		if receiverName == structName {
			return funcDecl, nil
		}
	}

	return nil, fmt.Errorf("method %s not found on struct %s", methodName, structName)
}

func (u UpdateSourceGo) maybeUpdateSource(filename string) (bool, error) {
	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		return false, fmt.Errorf("error reading file %s: %w", filename, err)
	}

	// Create a file set and parse the source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		return false, fmt.Errorf("error parsing file %s: %w", filename, err)
	}

	// Track modifications
	var parametersMethod, configureMethod *ast.FuncDecl
	var structName string
	var structEnd int

	// Inspect the AST to find methods
	ast.Inspect(file, func(n ast.Node) bool {
		if !u.implementsSource(file, n) {
			return true
		}

		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			panic(errors.New("node is not a TypeSpec"))
		}
		structName = typeSpec.Name.Name
		structEnd = fset.Position(typeSpec.End()).Offset

		parametersMethod, err = u.getMethod(file, n, "Parameters")
		configureMethod, err = u.getMethod(file, n, "Configure")

		return false
	})

	// If no modifications needed, skip this file
	if parametersMethod == nil && configureMethod == nil {
		return false, nil
	}

	// Create a new buffer to write modified source
	var newContent bytes.Buffer
	if err := format.Node(&newContent, fset, file); err != nil {
		return false, fmt.Errorf("error formatting AST: %w", err)
	}
	modifiedSource := newContent.String()

	// Remove Parameters method if it exists
	var parametersMethodStart, parametersMethodEnd int
	if parametersMethod != nil {
		parametersMethodStart = fset.Position(parametersMethod.Pos()).Offset
		parametersMethodEnd = fset.Position(parametersMethod.End()).Offset
		modifiedSource = strings.Replace(modifiedSource, modifiedSource[parametersMethodStart:parametersMethodEnd], "", 1)
	}

	// 1
	// 2

	// Add TODO comment for Configure method
	if configureMethod != nil {
		start := fset.Position(configureMethod.Pos()).Offset

		if parametersMethodStart != 0 && parametersMethodStart < start {
			start = start - (parametersMethodEnd - parametersMethodStart + 1)
		}

		todoComment := "\n// TODO: This method needs to be removed. If there's any custom logic in Configure(),\n" +
			"// it needs to be moved to the configuration struct in the Validate() method."
		modifiedSource = modifiedSource[:start] + todoComment + modifiedSource[start:]
	}

	// Add Config method before the last closing brace
	configMethodTemplate := fmt.Sprintf(
		`

func (s *%s) Config() sdk.SourceConfig {
	return s.config
}
`, structName)

	modifiedSource = modifiedSource[:structEnd] + configMethodTemplate + modifiedSource[structEnd:]

	// Write modified content back to file
	if err := os.WriteFile(filename, []byte(modifiedSource), 0644); err != nil {
		return false, fmt.Errorf("error writing modified file %s: %w", filename, err)
	}

	return true, nil
}
