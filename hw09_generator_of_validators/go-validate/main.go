package main

import (
	"github.com/fixme_my_friend/hw09_generator_of_validators/go-validate/generation"
	"github.com/fixme_my_friend/hw09_generator_of_validators/go-validate/inspections"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	path := "models/models.go"
	res := inspections.GetStructureMetadata(path)
	s, err := generation.GenerateValidationCode(res)
	if err != nil {
		log.Fatal(err)
	}

	path = strings.TrimSuffix(path, ".go")
	err = ioutil.WriteFile(path+"_validation_generated.go", []byte(s), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
