package fieldtype

import (
	"fmt"
	"reflect"

	"github.com/zl-leaf/gorm-dao/internal/method"
	"gorm.io/gorm/schema"
)

type baseField struct {
	field *schema.Field
}

// elemType 返回字段类型的元素类型（若为指针则解引用）
func (f *baseField) elemType() reflect.Type {
	typ := f.field.FieldType
	if typ.Kind() == reflect.Ptr {
		return typ.Elem()
	}
	return typ
}

// elemTypeName 返回用于参数的类型名（如 "string", "time.Time"）
func (f *baseField) elemTypeName() string {
	typ := f.elemType()
	if typ.PkgPath() != "" {
		return typ.String()
	}
	return typ.Name()
}

// paramType 返回单值参数的类型字符串（统一为元素类型，如 "int", "time.Time"，指针字段也为 "int" 而非 "*int"）
func (f *baseField) paramType() string {
	return f.elemTypeName()
}

// isPointer 返回字段是否为指针类型
func (f *baseField) isPointer() bool {
	return f.field.FieldType.Kind() == reflect.Ptr
}

func (f *baseField) eq() *method.Method {
	typeStr := f.paramType()
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: typeStr},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sEq", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Eq{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name),
	}
}

func (f *baseField) neq() *method.Method {
	typeStr := f.paramType()
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: typeStr},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sNeq", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Neq{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name),
	}
}

func (f *baseField) like() *method.Method {
	typeStr := f.paramType()
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: typeStr},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sLike", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Like{Column: \"%s\",Value: %s}", f.field.DBName, "\"%\"+"+params[0].Name+"+\"%\""),
	}
}

func (f *baseField) prefixLike() *method.Method {
	typeStr := f.paramType()
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: typeStr},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sPrefixLike", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Like{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name+"+\"%\""),
	}
}

func (f *baseField) notLike() *method.Method {
	typeStr := f.paramType()
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: typeStr},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sNotLike", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Not(clause.Like{Column: \"%s\",Value:  %s})", f.field.DBName, "\"%\"+"+params[0].Name+"+\"%\""),
	}
}

func (f *baseField) in() *method.Method {
	elemName := f.elemTypeName()
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: elemName, IsArray: true},
	}
	toSliceFnStr := fmt.Sprintf(`func(v []%s) []interface{} {
		ret := make([]interface{}, len(v))
		for i, item := range v {
			ret[i] = item
		}
		return ret
	}(%s)`, params[0].Type, params[0].Name)
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sIn", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.IN{Column: \"%s\",Values:  %s}", f.field.DBName, toSliceFnStr),
	}
}

func (f *baseField) notIn() *method.Method {
	elemName := f.elemTypeName()
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: elemName, IsArray: true},
	}
	toSliceFnStr := fmt.Sprintf(`func(v []%s) []interface{} {
		ret := make([]interface{}, len(v))
		for i, item := range v {
			ret[i] = item
		}
		return ret
	}(%s)`, params[0].Type, params[0].Name)
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sNotIn", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Not((clause.IN{Column: \"%s\",Values:  %s}))", f.field.DBName, toSliceFnStr),
	}
}

func (f *baseField) gt() *method.Method {
	typeStr := f.paramType()
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: typeStr},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sGt", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Gt{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name),
	}
}

func (f *baseField) gte() *method.Method {
	typeStr := f.paramType()
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: typeStr},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sGte", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Gte{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name),
	}
}

func (f *baseField) lt() *method.Method {
	typeStr := f.paramType()
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: typeStr},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sLt", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Lt{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name),
	}
}

func (f *baseField) lte() *method.Method {
	typeStr := f.paramType()
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: typeStr},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sLte", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Lte{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name),
	}
}

func (f *baseField) between() *method.Method {
	typeStr := f.paramType()
	params := method.Params{
		{Name: "left", Type: typeStr},
		{Name: "right", Type: typeStr},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sBetween", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Expr{SQL: \"%s Between ? AND ?\", Vars:[]interface{} {%s, %s}}", f.field.DBName, params[0].Name, params[1].Name),
	}
}

func (f *baseField) isNull() *method.Method {
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sIsNull", f.field.Name),
		Params:     method.Params{},
		Query:      fmt.Sprintf("clause.Expr{SQL: \"%s IS NULL\"}", f.field.DBName),
	}
}

func (f *baseField) isNotNull() *method.Method {
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sIsNotNull", f.field.Name),
		Params:     method.Params{},
		Query:      fmt.Sprintf("clause.Expr{SQL: \"%s IS NOT NULL\"}", f.field.DBName),
	}
}
