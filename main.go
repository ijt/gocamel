// The gocamel program converts snake_case to camelCase in Go programs.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	if err := gocamel(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func gocamel() error {
	willPrintAST := flag.Bool("print", false, "whether to print out the AST")
	flag.Parse()
	if flag.NArg() == 0 {
		return errors.New("usage: gocamel file1.go ... fileN.go")
	}

	for _, filename := range flag.Args() {
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("reading %s: %w", filename, err)
		}

		newContents, err := snakeCaseToCamelCaseFile(filename, contents, *willPrintAST)
		if err != nil {
			return fmt.Errorf("processing %s: %w", filename, err)
		}

		if err := ioutil.WriteFile(filename, newContents, 0640); err != nil {
			return fmt.Errorf("writing new version of %s: %w", filename, err)
		}
	}
	return nil
}

func snakeCaseToCamelCaseFile(filename string, contents []byte, willPrintAST bool) ([]byte, error) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	fileAST, err := parser.ParseFile(fset, filename, contents, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", filename, err)
	}

	if willPrintAST {
		if err := ast.Print(fset, fileAST); err != nil {
			return nil, fmt.Errorf("printing out AST: %w", err)
		}
	}

	ast.Inspect(fileAST, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			if strings.HasPrefix(x.Name, "Test") {
				return true
			}
			x.Name = snakeToCamel(x.Name)
		}
		return true
	})

	// Overwrite the original file from the modified AST.
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, fileAST); err != nil {
		return nil, fmt.Errorf("formatting AST for %s: %w", filename, err)
	}
	return buf.Bytes(), nil
}

var rx = regexp.MustCompile(`(\w)_(\w)`)

func snakeToCamel(ident string) string {
	for {
		ident2 := rx.ReplaceAllStringFunc(ident, func(s string) string {
			return s[0:1] + strings.ToUpper(s[2:3])
		})
		if ident2 == ident {
			return ident
		}
		ident = ident2
	}
}
