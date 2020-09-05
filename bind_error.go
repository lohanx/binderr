package binderr

import (
	"github.com/go-playground/validator/v10"
)

type BindErrors struct {
	validationErrors validator.ValidationErrors
	messages         map[string]map[string]error
	errors           []error
	esl              int                       //errors length
	esi              map[string]map[string]int //errors index
}

// errors gin build error
// messages are custom error messages
// var messages = map[string]map[string]error{
//     "email":{
//         "required":errors.New("error"),
//         "email":errors.New("error"),
//     },
//     "password":{
//         "required":errors.New("error"),
//         "len":errors.New("error"),
//     },
// }
func New(errors error, messages map[string]map[string]error) *BindErrors {
	be := new(BindErrors)
	be.init(errors, messages)
	be.parse()
	return be
}

func (be *BindErrors) init(errors error, messages map[string]map[string]error) {
	defer func() {
		if err := recover(); err != nil {
			be.errors = []error{errors}
			be.esl = 1
		}
	}()
	be.validationErrors = errors.(validator.ValidationErrors)
	be.messages = messages
}

func (be *BindErrors) parse() {
	if be.errors != nil {
		return
	} else {
		be.errors = make([]error, len(be.validationErrors))
	}
	be.esi = make(map[string]map[string]int)
	for i, err := range be.validationErrors {
		field := err.Field()
		tag := err.Tag()
		if _, ok := be.esi[field]; !ok {
			be.esi[field] = make(map[string]int)
		}
		be.esi[field][tag] = i
		be.errors[i] = be.messages[field][tag]
	}
	be.esl = len(be.errors)
}

// FirstError return the first error
func (be *BindErrors) FirstError() error {
	if be.esl <= 0 {
		return nil
	}
	return be.errors[0]
}

// LastError return the last error
func (be *BindErrors) LastError() error {
	if be.esl <= 0 {
		return nil
	}
	return be.errors[be.esl-1]
}

// GetError return error based on field and tag
func (be *BindErrors) GetTagError(field, tag string) error {
	if tags, ok := be.esi[field]; ok {
		if index, ok := tags[tag]; ok {
			return be.errors[index]
		}
	}
	return nil
}

// GetErrors return all errors of the field
func (be *BindErrors) GetFiledErrors(field string) []error {
	if tags, ok := be.esi[field]; ok {
		var errors = make([]error, 0, len(tags))
		for _, i := range tags {
			errors = append(errors, be.errors[i])
		}
		return errors
	}
	return nil
}

// Errors return all errors
func (be *BindErrors) Errors() []error {
	return be.errors
}

// Len return errors length
func (be *BindErrors) Len() int {
	return be.esl
}
