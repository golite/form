// Copyright 2015 Alex Browne and Soroush Pour.
// Allrights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package form

import "fmt"

type InputValidation struct {
	Errors    []error
	Input     *Input
	Form      *Form
	InputName string
}

type ValidationError struct {
	Input     *Input
	InputName string
	msg       string
}

func (valErr ValidationError) Error() string {
	return valErr.msg
}

func (form *Form) Validate(inputName string) *InputValidation {
	return &InputValidation{
		Input:     form.Inputs[inputName],
		Form:      form,
		InputName: inputName,
	}
}

func (val *InputValidation) AddError(format string, args ...interface{}) {
	err := fmt.Errorf(format, args...)
	val.Errors = append(val.Errors, err)
	valErr := &ValidationError{
		Input:     val.Input,
		InputName: val.InputName,
		msg:       err.Error(),
	}
	val.Form.Errors = append(val.Form.Errors, valErr)
}

func (val *InputValidation) Required() *InputValidation {
	return val.Requiredf("%s is required.", val.InputName)
}

func (val *InputValidation) Requiredf(format string, args ...interface{}) *InputValidation {
	if val.Input == nil || val.Input.RawValue == "" {
		val.AddError(format, args...)
	}
	return val
}

func lessFunc(limit int) func(int) bool {
	return func(value int) bool {
		return value < limit
	}
}

func lessOrEqualFunc(limit int) func(int) bool {
	return func(value int) bool {
		return value <= limit
	}
}

func greaterFunc(limit int) func(int) bool {
	return func(value int) bool {
		return value > limit
	}
}

func greaterOrEqualFunc(limit int) func(int) bool {
	return func(value int) bool {
		return value >= limit
	}
}

func (val *InputValidation) Less(limit int) *InputValidation {
	return val.Lessf(limit, "%s must be less than %d.", val.InputName, limit)
}

func (val *InputValidation) Lessf(limit int, format string, args ...interface{}) *InputValidation {
	return val.validateInt(lessFunc(limit), format, args...)
}

func (val *InputValidation) LessOrEqual(limit int) *InputValidation {
	return val.LessOrEqualf(limit, "%s must be less than or equal to %d.", val.InputName, limit)
}

func (val *InputValidation) LessOrEqualf(limit int, format string, args ...interface{}) *InputValidation {
	return val.validateInt(lessOrEqualFunc(limit), format, args...)
}

func (val *InputValidation) Greater(limit int) *InputValidation {
	return val.Greaterf(limit, "%s must be greater than %d.", val.InputName, limit)
}

func (val *InputValidation) Greaterf(limit int, format string, args ...interface{}) *InputValidation {
	return val.validateInt(greaterFunc(limit), format, args...)
}

func (val *InputValidation) GreaterOrEqual(limit int) *InputValidation {
	return val.GreaterOrEqualf(limit, "%s must be greater than or equal to %d.", val.InputName, limit)
}

func (val *InputValidation) GreaterOrEqualf(limit int, format string, args ...interface{}) *InputValidation {
	return val.validateInt(greaterOrEqualFunc(limit), format, args...)
}

func (val *InputValidation) validateInt(validateFunc func(value int) bool, format string, args ...interface{}) *InputValidation {
	// If the input does not exist or is empty, skip this validation.
	if val.Input == nil || val.Input.RawValue == "" {
		return val
	}
	// Attempt to convert the input value to an integer.
	intVal, err := val.Input.Int()
	if err != nil {
		val.AddError("%s must be an integer.", val.InputName)
		return val
	}
	// Call validateFunc and if it returns false, add the appropriate error.
	if !validateFunc(intVal) {
		val.AddError(format, args...)
	}
	return val
}