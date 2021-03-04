package validator

import (
	"fmt"
	"reflect"
)

type ValidatorFunc func(reflect.StructField, reflect.Value) []error

var (
	Validators []ValidatorFunc
)

func init() {
	Validators = make([]ValidatorFunc, reflect.UnsafePointer)
	// for i := reflect.Invalid; i < reflect.UnsafePointer; i++ {
	//
	// }

	Validators[reflect.String] = validateString
}

type ErrorType uint

const (
	Invalid ErrorType = iota
	Unsupported
	ValidationFailed
)

const (
	ErrorMessageInvalid     = "Incorrect Kind: %v, must be a reflect.Struct"
	ErrorMessageUnsupported = "Kind %v of Field %s is not yet supported"
	ErrorMessageFailed      = "Validation failed for field %s"
)

type ValidationError struct {
	ErrorType  ErrorType
	FieldName  string
	FieldKind  reflect.Kind
	FieldValue interface{}
	Message    string
}

func (ve *ValidationError) Error() string {
	switch ve.ErrorType {
	case Unsupported:
		return ve.errorUnsupported()
	case ValidationFailed:
		return ve.errorValidationFailed()
	}

	return "Unspecified validation error"
}

func (ve *ValidationError) errorUnsupported() string {
	return fmt.Sprintf(ErrorMessageUnsupported, ve.FieldKind, ve.FieldName)
}

func (ve *ValidationError) errorValidationFailed() string {
	return fmt.Sprintf(ErrorMessageFailed, ve.FieldName)
}

func Validate(i interface{}) []error {
	validationErrors := make([]error, 0)

	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Struct {
		return append(validationErrors, fmt.Errorf("Incorrect Kind: %v, must be a reflect.Struct", val.Kind()))
	}

	t := reflect.TypeOf(i)

	for i := 0; i < t.NumField(); i++ {
		typeField := t.Field(i)

		if isExportedField(typeField) {
			if valFunc := Validators[typeField.Type.Kind()]; valFunc != nil {
				validationErrors = append(validationErrors, (valFunc(typeField, val.Field(i)))...)
			} else {
				validationErrors = append(validationErrors, &ValidationError{
					FieldName: typeField.Name,
					ErrorType: Unsupported,
					FieldKind: typeField.Type.Kind(),
				})
			}
		}
	}
	return validationErrors
}

func isExportedField(field reflect.StructField) bool {
	return field.PkgPath == ""
}

func validateString(f reflect.StructField, v reflect.Value) []error {
	return nil
}
