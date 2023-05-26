package parser

import (
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func SourceFileAddImport(source_file string, out_file string) {
	fs := token.NewFileSet()

	file, err := os.ReadFile(source_file)
	if err != nil {
		log.Fatal("source file can't be read", err)
	}

	f, err := parser.ParseFile(fs, "", file, parser.AllErrors)
	if err != nil {
		log.Fatalln(err)
	}

	AddImport(f, "\"hook.com/hook/parser\"", "parser/ast.go")

	// printer.Fprint(os.Stdout, fs, f)
	out, err := os.Create(out_file)

	if err != nil {
		log.Fatalln("out_file can't be wrote ", err)
	}
	printer.Fprint(out, fs, f)
}
