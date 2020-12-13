package main

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	FirstName string `validate:"required"`
}

func main() {
	user := &User{FirstName: "firstName"}
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			fmt.Println(err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			fmt.Println(errs)
		}
	}
}
