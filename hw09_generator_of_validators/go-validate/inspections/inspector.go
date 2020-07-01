package inspections

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

func GetStructureMetadata(path string) []StructMetadata {
	file := getAstFromFile(path)
	astData := getAstStructureMetadata(file)

	var res []StructMetadata
	for _, metadata := range astData {
		s := StructMetadata{}
		s.Name = metadata.TypeSpec.Name.Name
		for _, field := range metadata.StructType.Fields.List {
			for _, name := range field.Names {
				if field.Tag != nil && strings.Contains(field.Tag.Value, "validate:") {
					var typeName string
					switch r := field.Type.(type) {
					case *ast.Ident:
						typeName = r.Name
					case *ast.ArrayType:
						typeName = "[]" + r.Elt.(*ast.Ident).Name
					}
					s.Fields = append(s.Fields, FieldMetadata{
						Name: name.Name,
						Type: typeName,
						Tag:  field.Tag.Value,
					})
				}
			}
		}
		res = append(res, s)
	}

	return res
}

func getAstFromFile(path string) *ast.File {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, path, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func getAstStructureMetadata(file *ast.File) []StructAstMetadata {
	needValidateStructures := make([]StructAstMetadata, 0)
	ast.Inspect(file, func(node ast.Node) bool {
		return checkNode(node, &needValidateStructures)
	})

	return needValidateStructures
}

func checkNode(n ast.Node, needValidateStructures *[]StructAstMetadata) bool {
	currType, ok := n.(*ast.TypeSpec)
	if !ok {
		return true
	}
	currStruct, ok := currType.Type.(*ast.StructType)
	if !ok {
		return true
	}

DONE:
	for _, field := range currStruct.Fields.List {
		if field.Tag == nil {
			continue
		}
		if strings.Contains(field.Tag.Value, "validate:") {
			*needValidateStructures = append(*needValidateStructures, StructAstMetadata{
				StructType: currStruct,
				TypeSpec:   currType,
			})
			break DONE
		}
	}
	return false
}

type StructAstMetadata struct {
	StructType *ast.StructType
	TypeSpec   *ast.TypeSpec
}
