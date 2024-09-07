package fieldtype

import (
	"github.com/zl-leaf/gorm-dao/internal/method"
	"gorm.io/gorm/schema"
)

type StringField struct {
	baseField
}

func NewStringField(field *schema.Field) *StringField {
	return &StringField{baseField{field: field}}
}

func (f *StringField) GetWhereMethods() []*method.Method {
	return []*method.Method{f.eq(), f.neq(), f.like(), f.prefixLike(), f.notLike(), f.in(), f.notIn()}
}
