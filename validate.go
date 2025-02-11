package main

import (
	"github.com/go-playground/validator/v10"
)

// Global singletons
var (
	validate *validator.Validate
)

func Setup() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// Validate returns the global validator instance.
func Validate() *validator.Validate {
	return validate
}
