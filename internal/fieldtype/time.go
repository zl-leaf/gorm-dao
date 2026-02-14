package fieldtype

import (
	"github.com/zl-leaf/gorm-dao/internal/method"
	"gorm.io/gorm/schema"
)

type TimeField struct {
	baseField
}

func NewTimeField(field *schema.Field) *TimeField {
	return &TimeField{baseField{field: field}}
}

func (f *TimeField) GetWhereMethods() []*method.Method {
	methods := []*method.Method{f.eq(), f.gt(), f.gte(), f.lt(), f.lte(), f.between()}
	if f.isPointer() {
		methods = append(methods, f.isNull(), f.isNotNull())
	}
	return methods
}
