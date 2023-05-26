package parser

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func SourceFileAdd(source_file string, out_file string) {
	fs := token.NewFileSet()

	file, err := os.ReadFile(source_file)
	if err != nil {
		log.Fatal("source file can't be read", err)
	}

	f, err := parser.ParseFile(fs, "", file, parser.AllErrors)
	if err != nil {
		log.Fatalln(err)
	}

	/* ----------------
	all changes on ast should be here:
		1.AddImport
	*/

	AddImport(f, "\"hook.com/hook/parser\"", "parser/ast.go")

	function_name := FindFunctions(f)
	fmt.Println(function_name)

	gofunction_name := findGoNames(f)
	fmt.Println(gofunction_name)

	printAst(f)

	// printer.Fprint(os.Stdout, fs, f)
	out, err := os.Create(out_file)

	if err != nil {
		log.Fatalln("out_file can't be wrote ", err)
	}
	printer.Fprint(out, fs, f)
}
