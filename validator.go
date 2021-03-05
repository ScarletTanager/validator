package validator

import (
	"fmt"
	// "os"
	"reflect"
	"strconv"
	"strings"
)

type ValidatorFunc func(reflect.StructField, reflect.Value) []*ValidationError

var (
	Validators []ValidatorFunc
)

func init() {
	Validators = make([]ValidatorFunc, reflect.UnsafePointer)

	Validators[reflect.String] = validateString
	Validators[reflect.Int] = validateInt
}

type ErrorType uint

const (
	Invalid ErrorType = iota
	Unsupported
	ValidationFailed
)

const (
	ErrorMessageInvalid          = "Incorrect Kind: %v, must be a reflect.Struct"
	ErrorMessageUnsupported      = "Kind %v of Field %s is not yet supported"
	ErrorMessageValidationFailed = "Validation failed for field %s"
)

type ValidationError struct {
	ErrorType ErrorType
	Message   string
}

func (ve *ValidationError) Error() string {
	return ve.Message
}

func Validate(i interface{}) []*ValidationError {
	validationErrors := make([]*ValidationError, 0)

	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Struct {
		return append(validationErrors, &ValidationError{
			ErrorType: Invalid,
			Message:   fmt.Sprintf(ErrorMessageInvalid, val.Kind()),
		})
	}

	t := reflect.TypeOf(i)

	for i := 0; i < t.NumField(); i++ {
		typeField := t.Field(i)

		if isExportedField(typeField) && isValidatedField(typeField) {
			if valFunc := Validators[typeField.Type.Kind()]; valFunc != nil {
				validationErrors = append(validationErrors, (valFunc(typeField, val.Field(i)))...)
			} else {
				validationErrors = append(validationErrors, &ValidationError{
					ErrorType: Unsupported,
					Message:   fmt.Sprintf(ErrorMessageUnsupported, typeField.Type.Kind(), typeField.Name),
				})
			}
		}
	}
	return validationErrors
}

func isValidatedField(field reflect.StructField) bool {
	return field.Tag.Get("validator") != ""
}

func isExportedField(field reflect.StructField) bool {
	return field.PkgPath == ""
}

func validateString(f reflect.StructField, v reflect.Value) []*ValidationError {
	errs := make([]*ValidationError, 0)

	reqs := parseRequirements(f.Tag.Get("validator"))
	if reqs.required() && !reqs.allowEmpty() {
		if fmt.Sprintf("%v", v.Interface()) == "" {
			errs = append(errs, &ValidationError{
				ErrorType: ValidationFailed,
				Message:   fmt.Sprintf(ErrorMessageValidationFailed, f.Name),
			})
		}
	}
	return errs
}

func validateInt(f reflect.StructField, v reflect.Value) []*ValidationError {
	errs := make([]*ValidationError, 0)

	reqs := parseRequirements(f.Tag.Get("validator"))
	if reqs.required() {
		if gt, boundary := reqs.greaterThanInt(); gt {
			if i := (v.Interface()).(int); i <= boundary {
				errs = append(errs, &ValidationError{
					ErrorType: ValidationFailed,
					Message:   fmt.Sprintf(ErrorMessageValidationFailed, f.Name),
				})
			}
		}
	}
	return errs
}

func parseRequirements(tag string) Requirements {
	tagFields := strings.Split(tag, ",")
	if len(tagFields) == 1 && tagFields[0] == "" {
		return nil
	}

	reqs := make(Requirements)

	for i := 0; i < len(tagFields); i++ {
		f := tagFields[i]
		if f == "greaterthan" || f == "lessthan" {
			i++
			reqs.addRequirementWithParam(f, tagFields[i])
		} else {
			reqs.addRequirement(f)
		}
	}

	return reqs
}

type Requirements map[string]interface{}

// General requirements
func (r Requirements) required() bool {
	_, required := r["required"]
	return required
}

func (r Requirements) allowEmpty() bool {
	_, allowEmpty := r["allowempty"]
	return allowEmpty
}

// Specific to integer fields
func (r Requirements) greaterThanInt() (bool, int) {
	boundaryStr, greaterthan := r["greaterthan"]
	boundary, err := strconv.Atoi(boundaryStr.(string))
	if err != nil {
		return false, 0
	}
	return greaterthan, boundary
}

func (r Requirements) addRequirement(req string) {
	r.addRequirementWithParam(req, nil)
}

func (r Requirements) addRequirementWithParam(req string, param interface{}) {
	r[req] = param
}
