package fieldtype

import (
	"errors"
	"reflect"

	"github.com/zl-leaf/gorm-dao/internal/method"
	"gorm.io/gorm/schema"
)

type IConditionField interface {
	GetWhereMethods() []*method.Method
}

// elemType 返回字段类型的元素类型（若为指针则解引用）
func elemType(field *schema.Field) reflect.Type {
	typ := field.FieldType
	if typ.Kind() == reflect.Ptr {
		return typ.Elem()
	}
	return typ
}

func GetConditionFieldBySchemaField(field *schema.Field) (IConditionField, error) {
	typ := elemType(field)
	typeKey := typ.Name()
	if typ.PkgPath() != "" {
		typeKey = typ.String() // 如 "time.Time"
	}
	switch typeKey {
	case "string":
		return NewStringField(field), nil
	case "int", "int32", "int64", "float", "float32", "float64", "uint", "uint32", "uint64":
		return NewNumberField(field), nil
	case "time.Time":
		return NewTimeField(field), nil
	}
	return nil, errors.New("无合适的field")
}
