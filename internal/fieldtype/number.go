package fieldtype

import (
	"github.com/zl-leaf/gorm-dao/internal/method"
	"gorm.io/gorm/schema"
)

type NumberField struct {
	baseField
}

func NewNumberField(field *schema.Field) *NumberField {
	return &NumberField{baseField{field: field}}
}

func (f *NumberField) GetWhereMethods() []*method.Method {
	return []*method.Method{f.eq(), f.neq(), f.in(), f.notIn(), f.lt(), f.lte(), f.gt(), f.gte(), f.between()}
}
