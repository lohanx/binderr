package binderr

import "github.com/go-playground/validator/v10"

type BindErrors struct {
	ValidationErrors validator.ValidationErrors
	ErrorMsgs        map[string]map[string]error
	errors           []error
	errorsLen        int
	errorIndex       map[string]map[string]int
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
	be := &BindErrors{}
	be.apply(errors, messages)
	be.init()
	return be
}

func (be *BindErrors) apply(errors error, messages map[string]map[string]error) {
	defer func() {
		if err := recover(); err != nil {
			be.errors = make([]error, 0, 1)
			be.errors = append(be.errors, errors)
			be.errorsLen = 1
		}
	}()
	be.ValidationErrors = errors.(validator.ValidationErrors)
	be.ErrorMsgs = messages
}

func (be *BindErrors) init() {
	if be.errors != nil {
		return
	} else {
		be.errors = make([]error, 0, len(be.ValidationErrors))
	}
	if be.errorIndex == nil {
		be.errorIndex = make(map[string]map[string]int)
	}
	for i, err := range be.ValidationErrors {
		field := err.Field()
		tag := err.Tag()
		if _, ok := be.errorIndex[field]; !ok {
			be.errorIndex[field] = make(map[string]int)
		}
		be.errorIndex[field][tag] = i
		be.errors = append(be.errors, be.ErrorMsgs[field][tag])
	}
	be.errorsLen = len(be.errors)
}

// LastError return the last error
func (be *BindErrors) LastError() error {
	if be.errorsLen <= 0 {
		return nil
	}
	return be.errors[be.errorsLen-1]
}

// GetError return error based on field and tag
func (be *BindErrors) GetError(field, tag string) error {
	if errors, ok := be.errorIndex[field]; ok {
		if index, ok := errors[tag]; ok {
			return be.errors[index]
		}
	}
	return nil
}

// GetErrors return all errors of the field
func (be *BindErrors) GetErrors(field string) []error {
	if errsIndex, ok := be.errorIndex[field]; ok {
		var errors = make([]error, 0, len(errsIndex))
		for _, index := range errsIndex {
			errors = append(errors, be.errors[index])
		}
		return errors
	}
	return nil
}

// FirstError return the first error
func (be *BindErrors) FirstError() error {
	if be.errorsLen <= 0 {
		return nil
	}
	return be.errors[0]
}

// Errors return all errors
func (be *BindErrors) Errors() []error {
	return be.errors
}
