package models

import (
	"regexp"
	"strconv"
)

type ValidationError struct {
	Field    string
	Validate string
}

func (u User) Validate() (validationErrors []ValidationError, err error) {
	if len(u.ID) == 36 {
		validationErrors = append(validationErrors, ValidationError{
			Field:    "ID",
			Validate: "len",
		})
	}
	if u.Age < 18 {
		validationErrors = append(validationErrors, ValidationError{
			Field:    "Age",
			Validate: "min",
		})
	}
	if u.Age > 50 {
		validationErrors = append(validationErrors, ValidationError{
			Field:    "Age",
			Validate: "max",
		})
	}
	if matched, err := regexp.MatchString("^\\w+@\\w+\\.\\w+$", u.Email); err != nil || !matched {
		validationErrors = append(validationErrors, ValidationError{
			Field:    "Email",
			Validate: "regexp",
		})
	}
	for _, elem := range u.Phones {
		if len(elem) == 11 {
			validationErrors = append(validationErrors, ValidationError{
				Field:    "Phones",
				Validate: "len",
			})
		}
	}
	return
}

func (a App) Validate() (validationErrors []ValidationError, err error) {
	if len(a.Version) == 5 {
		validationErrors = append(validationErrors, ValidationError{
			Field:    "Version",
			Validate: "len",
		})
	}
	return
}

func (r Response) Validate() (validationErrors []ValidationError, err error) {
	if _, ok := map[string]struct{}{"200": {}, "404": {}, "500": {}}[strconv.Itoa(r.Code)]; !ok {
		validationErrors = append(validationErrors, ValidationError{
			Field:    "Code",
			Validate: "in",
		})
	}
	return
}
