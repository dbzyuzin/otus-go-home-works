package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
)

var vf = `
func bla bla bla {
 {{.}}
}
`

var templ = `
{{define "validate field"}}

{{end}}
{{define "validate slice"}}

{{end}}


text

{{template "validate func" .}}
{{template "validate func" "text"}}
`

func main() {

	tmpl, err := parseTemplate()
	if err != nil {
		log.Fatal(err)
	}

	//file, err := os.Create("res.go")
	//if err != nil {
	//	log.Fatal(err)
	//}
	data := []ValidationMethod{
		{StructName: "App",
			Fields: []ValidationField{{Steps: []ValidationCondition{{IsSlice: false, Condition: `len(s.Naame) > 5`, FieldName: "Naame", ErrorName: "ErrStringLen"}}}}}}

	fmt.Printf("$+v", data)
	println()
	var builder bytes.Buffer
	err = tmpl.ExecuteTemplate(&builder, "main", FileTemplateMetadata{Package: "models", Structs: data})
	if err != nil {
		log.Fatalf("text %w", err)
	}
	s, err := format.Source(builder.Bytes())
	if err != nil {
		log.Fatalf("source err %w", err)
	}
	err = ioutil.WriteFile("res.go", s, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
