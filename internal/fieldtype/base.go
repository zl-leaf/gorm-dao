package fieldtype

import (
	"fmt"

	"github.com/zl-leaf/gorm-dao/internal/method"
	"gorm.io/gorm/schema"
)

type baseField struct {
	field *schema.Field
}

func (f *baseField) eq() *method.Method {
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: f.field.FieldType.Name()},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sEq", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Eq{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name),
	}
}

func (f *baseField) neq() *method.Method {
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: f.field.FieldType.Name()},
	}

	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sNeq", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Neq{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name),
	}
}

func (f *baseField) like() *method.Method {
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: f.field.FieldType.Name()},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sLike", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Like{Column: \"%s\",Value: %s}", f.field.DBName, "\"%\"+"+params[0].Name+"+\"%\""),
	}
}

func (f *baseField) prefixLike() *method.Method {
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: f.field.FieldType.Name()},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sPrefixLike", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Like{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name+"+\"%\""),
	}
}

func (f *baseField) notLike() *method.Method {
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: f.field.FieldType.Name()},
	}
	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sNotLike", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Not(clause.Like{Column: \"%s\",Value:  %s})", f.field.DBName, "\"%\"+"+params[0].Name+"+\"%\""),
	}
}

func (f *baseField) in() *method.Method {
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: f.field.FieldType.Name(), IsArray: true},
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
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: f.field.FieldType.Name(), IsArray: true},
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
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: f.field.FieldType.Name()},
	}

	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sGt", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Gt{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name),
	}
}

func (f *baseField) gte() *method.Method {
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: f.field.FieldType.Name()},
	}

	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sGte", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Gte{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name),
	}
}

func (f *baseField) lt() *method.Method {
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: f.field.FieldType.Name()},
	}

	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sLt", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Lt{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name),
	}
}

func (f *baseField) lte() *method.Method {
	params := method.Params{
		{Name: method.ToArgsName(f.field.Name), Type: f.field.FieldType.Name()},
	}

	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sLte", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Lte{Column: \"%s\",Value:  %s}", f.field.DBName, params[0].Name),
	}
}

func (f *baseField) between() *method.Method {
	params := method.Params{
		{Name: "left", Type: f.field.FieldType.Name()},
		{Name: "right", Type: f.field.FieldType.Name()},
	}

	return &method.Method{
		Field:      f.field,
		MethodName: fmt.Sprintf("Where%sBetween", f.field.Name),
		Params:     params,
		Query:      fmt.Sprintf("clause.Expr{SQL: \"%s Between ? AND ?\", Vars:[]interface{} {%s, %s}}", f.field.DBName, params[0].Name, params[1].Name),
	}
}
