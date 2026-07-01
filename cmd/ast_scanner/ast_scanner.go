package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Export struct {
	Module string
	Pkg    string
	Name   string
	Type   string // struct, interface
}

type Usage struct {
	FromModule string
	FromPkg    string
	ToModule   string
	ToPkg      string
	Name       string
	File       string
	Line       int
}

func main() {
	fset := token.NewFileSet()
	
	exports := make(map[string]map[string]Export) // pkg -> name -> Export
	
	// first pass: find exports
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}
		if !strings.HasPrefix(path, "sovereign-") {
			return nil
		}
		if strings.Contains(path, "vendor") || strings.Contains(path, ".git") {
			return filepath.SkipDir
		}
		
		pkgs, err := parser.ParseDir(fset, path, nil, 0)
		if err != nil {
			return nil
		}
		
		modName := strings.Split(path, string(os.PathSeparator))[0]
		
		for pkgName, pkg := range pkgs {
			if strings.HasSuffix(pkgName, "_test") { continue }
			if exports[pkgName] == nil {
				exports[pkgName] = make(map[string]Export)
			}
			for _, file := range pkg.Files {
				ast.Inspect(file, func(n ast.Node) bool {
					gen, ok := n.(*ast.GenDecl)
					if ok && gen.Tok == token.TYPE {
						for _, spec := range gen.Specs {
							ts, ok := spec.(*ast.TypeSpec)
							if ok && ts.Name.IsExported() {
								var typ string
								switch ts.Type.(type) {
								case *ast.StructType:
									typ = "struct"
								case *ast.InterfaceType:
									typ = "interface"
								}
								if typ != "" {
									exports[pkgName][ts.Name.Name] = Export{
										Module: modName,
										Pkg: pkgName,
										Name: ts.Name.Name,
										Type: typ,
									}
								}
							}
						}
					}
					return true
				})
			}
		}
		return nil
	})

	var usages []Usage

	// second pass: find usages
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if !strings.HasPrefix(path, "sovereign-") {
			return nil
		}
		if strings.Contains(path, "vendor") || strings.Contains(path, ".git") {
			return nil
		}
		
		modName := strings.Split(path, string(os.PathSeparator))[0]
		
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil
		}
		
		pkgName := file.Name.Name
		
		ast.Inspect(file, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if ok {
				ident, ok := sel.X.(*ast.Ident)
				if ok {
					if exps, found := exports[ident.Name]; found {
						if exp, foundName := exps[sel.Sel.Name]; foundName {
							if exp.Module != modName {
								usages = append(usages, Usage{
									FromModule: modName,
									FromPkg: pkgName,
									ToModule: exp.Module,
									ToPkg: exp.Pkg,
									Name: exp.Name,
									File: path,
									Line: fset.Position(sel.Pos()).Line,
								})
							}
						}
					}
				}
			}
			return true
		})
		return nil
	})
	
	// Create report
	report, err := os.Create("CONTRACT_BREAKAGE_REPORT.md")
	if err != nil {
		panic(err)
	}
	defer report.Close()
	
	fmt.Fprintln(report, "# Contract Breakage Report")
	fmt.Fprintln(report, "## Phase III: Schema and Contract Validation")
	fmt.Fprintln(report, "\n### Cross-Module Struct/Interface Usages Detected:")
	
	for _, u := range usages {
		fmt.Fprintf(report, "- **%s** (%s:%d) uses `%s.%s` from **%s**\n", u.FromModule, u.File, u.Line, u.ToPkg, u.Name, u.ToModule)
	}
	
	fmt.Fprintln(report, "\n### Analysis")
	if len(usages) > 0 {
		fmt.Fprintln(report, "Contract interfaces have been mapped. The above structs and interfaces cross module boundaries. Any change to their exported fields will result in a contract breakage.")
	} else {
		fmt.Fprintln(report, "No cross-module struct/interface usage found. Modules are perfectly decoupled or use generic interfaces (e.g., protobuf/json).")
	}
	fmt.Println("Report generated successfully.")
}
