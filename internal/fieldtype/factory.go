package fieldtype

import (
	"errors"

	"github.com/zl-leaf/gorm-dao/internal/method"
	"gorm.io/gorm/schema"
)

type IConditionField interface {
	GetWhereMethods() []*method.Method
}

func GetConditionFieldBySchemaField(field *schema.Field) (IConditionField, error) {
	switch field.FieldType.Name() {
	case "string":
		return NewStringField(field), nil
	case "int", "int32", "int64", "float", "float32", "float64", "uint", "uint32", "uint64":
		return NewNumberField(field), nil
	}
	return nil, errors.New("无合适的field")
}
