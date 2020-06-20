package main

import (
	"github.com/fixme_my_friend/hw09_generator_of_validators/go-validate/inspections"
)

func main() {
	_ = inspections.GetStructureMetadata("models/models.go")

	//fmt.Println()
	//for _, structure := range needValidateStructures {
	//	//fmt.Printf("func (s *%s) Validate() ([]ValidationError, error) {\n", structure.spec.Name)
	//	for _, field := range structure.s.Fields.List {
	//		fmt.Printf("names %v", field.Names)
	//		switch val := field.Type.(type) {
	//		case *ast.Ident:
	//			fmt.Printf("ident %v\n", val.Name)
	//		case *ast.ArrayType:
	//			fmt.Printf("array %v\n", val.Elt.(*ast.Ident).Name)
	//		}
	//		println()
	//		//fmt.Printf("%v\n", field.Type)
	//	}
	//}
}
