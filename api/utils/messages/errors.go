package messages

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

func ErrorsResponse(err error) interface{} {
	invalidSyntax := ""
	if strings.Index(err.Error(), "EOF") > -1 {
		invalidSyntax = "JSON is Empty"
		return invalidSyntax
	}

	if strings.Index(err.Error(), "unmarshal") > -1 {
		invalidSyntax = "Invalid Unmarshal"
		return invalidSyntax
	}

	if _, ok := err.(*json.SyntaxError); ok {
		invalidSyntax = "Invalid Syntax"
		return invalidSyntax
	}

	ve := err.(validator.ValidationErrors)
	invalidFields := make([]map[string]string, 0)
	message := ""
	for _, v := range ve {
		errors := map[string]string{}
		if v.ActualTag() == "required" {
			message = fmt.Sprintf("%v is required", strcase.ToLowerCamel(v.Field()))
		} else {
			message = fmt.Sprintf("%v has to be %v", strcase.ToLowerCamel(v.Field()), v.ActualTag())
		}

		errors[strcase.ToLowerCamel(v.Field())] = message
		invalidFields = append(invalidFields, errors)
	}

	return invalidFields
}
