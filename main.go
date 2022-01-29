// The gocamel program converts snake_case to camelCase in Go programs.
package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	if err := gocamel(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func gocamel() error {
	willPrint := flag.Bool("print", false, "whether to print out the AST")
	flag.Parse()
	if flag.NArg() == 0 {
		return errors.New("usage: gocamel file1.go ... fileN.go")
	}

	for _, filename := range flag.Args() {
		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("opening %s: %w", filename, err)
		}
		// Create the AST by parsing src.
		fset := token.NewFileSet() // positions are relative to fset
		fileAST, err := parser.ParseFile(fset, filename, f, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", filename, err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("closing %s: %w", filename, err)
		}

		if *willPrint {
			if err := ast.Print(fset, fileAST); err != nil {
				return fmt.Errorf("printing out AST: %w", err)
			}
		}

		// Overwrite the original file from the modified AST.
		f, err = os.Create(filename)
		if err != nil {
			return fmt.Errorf("opening %s for writing: %w", err)
		}
		if err := format.Node(f, fset, fileAST); err != nil {
			return fmt.Errorf("formatting AST for %s: %w", filename, err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("closing %s: %w", filename, err)
		}
	}
	return nil
}
