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
)

func main() {
	if err := gocamel(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func gocamel() error {
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

		// Create an ast.CommentMap from the ast.File's comments.
		// This helps keeping the association between comments
		// and AST nodes.
		cmap := ast.NewCommentMap(fset, fileAST, fileAST.Comments)

		// Remove the first variable declaration from the list of declarations.
		for i, decl := range fileAST.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.VAR {
				copy(fileAST.Decls[i:], fileAST.Decls[i+1:])
				fileAST.Decls = fileAST.Decls[:len(fileAST.Decls)-1]
				break
			}
		}

		// Use the comment map to filter comments that don't belong anymore
		// (the comments associated with the variable declaration), and create
		// the new comments list.
		fileAST.Comments = cmap.Filter(fileAST).Comments()

		// Overwrite the original file from the modified AST.
		var buf bytes.Buffer
		if err := format.Node(&buf, fset, fileAST); err != nil {
			return fmt.Errorf("formatting AST for %s: %w", filename, err)
		}
		if err := ioutil.WriteFile(filename, buf.Bytes(), 0640); err != nil {
			return fmt.Errorf("writing %s: %w", filename, err)
		}
	}
	return nil
}
