package generation

import (
	"errors"
	"fmt"
	"github.com/fixme_my_friend/hw09_generator_of_validators/go-validate/inspections"
	"strings"
)

var ErrTagCorrupted = errors.New("validate tag is corrupted")

type Validator struct {
	Type             string
	Name             string
	Value            string
	GenerateCode     func(string) string
	ValidationErrors []string
}

func GenerateValidationCode(structs []inspections.StructMetadata) (string, error) {
	var result strings.Builder

	result.WriteString("package models\n\nimport \"strconv\"\n\n")
	result.WriteString(`type ValidationError struct {
	Field string
	Validate string
}

`)

	for _, metadata := range structs {
		baseVarName := strings.ToLower(string([]rune(metadata.Name)[0]))
		result.WriteString(fmt.Sprintf("func (%s %s) Validate() (validationErrors []ValidationError, err error)  {", baseVarName, metadata.Name))

		for _, field := range metadata.Fields {
			switch field.Type {
			case "int":
				r, err := intValidationCode(field.Tag, field.Name, baseVarName+"."+field.Name)
				if err != nil {
					return "", err
				}
				for _, s := range r {
					result.WriteString("\n\t" + s)
				}
			case "string":
				r, err := stringValidationCode(field.Tag, field.Name, baseVarName+"."+field.Name)
				if err != nil {
					return "", err
				}
				for _, s := range r {
					result.WriteString("\n\t" + s)
				}
			case "[]int":
				r, err := intValidationCode(field.Tag, field.Name, "elem")
				if err != nil {
					return "", err
				}
				result.WriteString("\n\tfor _, elem := range " + baseVarName + "." + field.Name + " {")
				for _, s := range r {
					result.WriteString("\n\t\t" + s)
				}
				result.WriteString("\n\t}")
			case "[]string":
				r, err := stringValidationCode(field.Tag, field.Name, "elem")
				if err != nil {
					return "", err
				}
				result.WriteString("\n\tfor _, elem := range " + baseVarName + "." + field.Name + " {")
				for _, s := range r {
					result.WriteString("\n\t\t" + s)
				}
				result.WriteString("\n\t}")
			}
		}

		result.WriteString("\n\treturn\n}\n\n")
	}

	return result.String(), nil
}

func getValidateTagValue(tag string) string {
	tags := strings.Split(tag, " ")
	var validateTag string
	for _, s := range tags {
		if strings.Contains(s, "validate:") {
			validateTag = s
			break
		}
	}
	validateTag = strings.Trim(validateTag, "`")
	validateTag = strings.TrimPrefix(validateTag, "validate:\"")
	validateTag = strings.TrimSuffix(validateTag, "\"")
	return validateTag
}

func intValidationCode(tag string, fieldName string, fullFieldName string) ([]string, error) {
	template := "if %s {\n" +
		"\tvalidationErrors = append(validationErrors, ValidationError {\n" +
		"\t\tField:    \"%s\",\n" +
		"\t\tValidate: \"%s\",\n" +
		"\t})\n}"

	tag = getValidateTagValue(tag)
	var lines []string
	vSteps := strings.Split(tag, "|")
	for _, step := range vSteps {
		v := strings.Split(step, ":")
		if len(v) != 2 {
			return nil, ErrTagCorrupted
		}
		key, value := v[0], v[1]
		switch key {
		case "min":
			condition := fmt.Sprintf("%s < %s", fullFieldName, value)
			res := fmt.Sprintf(template, condition, fieldName, "min")
			lines = append(lines, strings.Split(res, "\n")...)
		case "max":
			condition := fmt.Sprintf("%s > %s", fullFieldName, value)
			res := fmt.Sprintf(template, condition, fieldName, "max")
			lines = append(lines, strings.Split(res, "\n")...)
		case "in":
			values := strings.Split(value, ",")
			m := make(map[string]struct{}, len(values))
			for _, s := range values {
				m[s] = struct{}{}
			}
			condition := fmt.Sprintf("_, ok := %#v[strconv.Itoa(%s)]; !ok", m, fullFieldName)
			res := fmt.Sprintf(template, condition, fieldName, "in")
			lines = append(lines, strings.Split(res, "\n")...)
		}
	}
	return lines, nil
}

func stringValidationCode(tag string, fieldName string, fullFieldName string) ([]string, error) {
	template := "if %s {\n" +
		"\tvalidationErrors = append(validationErrors, ValidationError {\n" +
		"\t\tField:    \"%s\",\n" +
		"\t\tValidate: \"%s\",\n" +
		"\t})\n}"

	tag = getValidateTagValue(tag)
	var lines []string
	vSteps := strings.Split(tag, "|")
	for _, step := range vSteps {
		v := strings.Split(step, ":")
		if len(v) != 2 {
			return nil, ErrTagCorrupted
		}
		key, value := v[0], v[1]
		switch key {
		case "len":
			condition := fmt.Sprintf("len(%s) == %s", fullFieldName, value)
			res := fmt.Sprintf(template, condition, fieldName, "len")
			lines = append(lines, strings.Split(res, "\n")...)
		case "regexp":
			condition := fmt.Sprintf("matched, err := regexp.MatchString(\"%s\", %s); err != nil || !matched", value, fullFieldName)
			res := fmt.Sprintf(template, condition, fieldName, "regexp")
			lines = append(lines, strings.Split(res, "\n")...)
		case "in":
			values := strings.Split(value, ",")
			m := make(map[string]struct{}, len(values))
			for _, s := range values {
				m[s] = struct{}{}
			}
			condition := fmt.Sprintf("_, ok := %#v[%s]; !ok", m, fullFieldName)
			res := fmt.Sprintf(template, condition, fieldName, "in")
			lines = append(lines, strings.Split(res, "\n")...)
		}
	}
	return lines, nil
}
