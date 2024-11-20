package validation

import (
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func ValidateStruct() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})

	return validate
}
