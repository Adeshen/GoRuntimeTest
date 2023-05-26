package parser

import (
	// _ "hook.com/hook/parser"
	"fmt"
	"go/ast"
	"go/token"
)

func AddImport(file *ast.File, packageName string, filename string) {
	print("AddImport ", filename, "\n")
	noImport := true
	toInsert := &ast.ImportSpec{
		Name: &ast.Ident{
			Name: "_",
		},
		Path: &ast.BasicLit{
			ValuePos: 0,
			Kind:     token.STRING,
			Value:    packageName,
		},
		EndPos: 0,
	}
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.GenDecl)
		if ok && fd.Tok == token.IMPORT {
			imports := make([]ast.Spec, 0, len(fd.Specs)+1)
			imports = append(imports, toInsert)
			imports = append(imports, fd.Specs...)
			fd.Specs = imports
			noImport = false
			fmt.Println("import true")
		}
	}
	if noImport {
		decls := make([]ast.Decl, 0, len(file.Decls)+1)
		imports := make([]ast.Spec, 0, 1)
		imports = append(imports, toInsert)
		decl := &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: imports,
		}
		decls = append(decls, decl)
		decls = append(decls, file.Decls...)
		file.Decls = decls
		fmt.Println("noImport")
	}
}

type FunctionVisitor struct {
	Functions []string
}

func (v *FunctionVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		v.Functions = append(v.Functions, n.Name.Name)
	}
	return v
}

func FindFunctions(node ast.Node) []string {

	// Create an instance of the FunctionVisitor
	visitor := &FunctionVisitor{}

	// Visit the AST nodes to find functions
	ast.Walk(visitor, node)

	// Return the list of functions
	return visitor.Functions
}
