package method

import (
	"strings"

	"gorm.io/gorm/schema"
)

type Method struct {
	Field *schema.Field

	MethodName string
	Params     Params
	Query      string
}

type Params []*Param
type Param struct { // (user model.User)
	PkgPath   string // package's path: internal/model
	Package   string // package's name: model
	Name      string // param's name: user
	Type      string // param's type: User
	IsArray   bool   // is array or not
	IsPointer bool   // is pointer or not
}

func (params Params) InputArgs() string {
	paramBody := make([]string, 0, len(params))

	for _, param := range params {
		paramName := param.Name
		paramType := param.Type

		if param.Package != "" {
			paramName = param.Package + "." + paramName
		}
		if param.IsPointer {
			paramName = "*" + paramName
		}
		if param.IsArray {
			paramType = "[]" + paramType
		}

		paramBody = append(paramBody, paramName+" "+paramType)
	}

	return strings.Join(paramBody, ",")
}

func (params Params) QueryArgs() string {
	paramNames := make([]string, len(params))
	for i, param := range params {
		paramNames[i] = param.Name
	}
	return strings.Join(paramNames, ",")
}

func ToArgsName(fieldName string) string {
	if fieldName == "ID" {
		return fieldName
	}
	return strings.ToLower(fieldName[:1]) + fieldName[1:]
}
