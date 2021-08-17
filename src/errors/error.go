package errors

import "fmt"

// Argument a key-value pair that is used for providing additional dynamic parameters.
// E.g. errArgs := []*app_err.Argument{{Key: "language", Value: "en"}}
// app_err.NewKnownf("MISSING_SUBJECT", "Subject is missing for language [%s]", errArgs)
type Argument struct {
	Key   string
	Value interface{}
}

// Known helps to handle expected logical errors.
// E.g. a client of the API doesn't provide required data, the Known error
// is being propagated to request's handler and returned as a clear message
// with a specific error code.
type Known struct {
	code    int
	message string
	args    map[string]interface{}
}

// NewKnown creates Known error
func NewKnown(code int, msg string) Known {
	return Known{
		code:    code,
		message: msg,
	}
}

// NewKnownf creates Known error with arguments
func NewKnownf(code int, format string, a []*Argument) Known {
	args := make(map[string]interface{})
	var as []interface{}
	for _, arg := range a {
		as = append(as, arg.Value)
		args[arg.Key] = arg.Value
	}
	msg := fmt.Sprintf(format, as...)

	return Known{
		code:    code,
		message: msg,
		args:    args,
	}
}

// Error returns error message string
func (e Known) Error() string {
	return e.message
}

// Code returns error code
func (e Known) Code() int {
	return e.code
}

// Args returns error message arguments
func (e Known) Args() map[string]interface{} {
	return e.args
}
