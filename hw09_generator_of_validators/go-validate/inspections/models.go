package inspections

type StructMetadata struct {
	Name   string
	Fields []FieldMetadata
}

type FieldMetadata struct {
	Name string
	Type string
	Tag  string `validate:"df"`
}
