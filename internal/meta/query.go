package meta

import (
	"fmt"
	"strings"

	"github.com/zl-leaf/gorm-dao/internal/fieldtype"
	"github.com/zl-leaf/gorm-dao/internal/method"
	"gorm.io/gorm/schema"
)

type QueryStruct struct {
	Schema *schema.Schema

	QueryPkgName string // 包名
	Name         string // Model对象名称
	StructName   string // 小写开头的Model名称
	WhereFns     []*method.Method
	OrderFns     []*method.Method
	PreloadFns   []*method.Method
}

// ParseQueryStruct 解析Query对象
func ParseQueryStruct(QueryPkgName string, s *schema.Schema) (*QueryStruct, error) {
	meta := &QueryStruct{
		Schema:       s,
		QueryPkgName: QueryPkgName,
		Name:         s.Name,
		StructName:   strings.ToLower(s.Name[:1]) + s.Name[1:],
		WhereFns:     make([]*method.Method, 0, 10),
		OrderFns:     make([]*method.Method, 0, 10),
		PreloadFns:   make([]*method.Method, 0, 10),
	}

	for _, field := range s.Fields {
		meta.WhereFns = append(meta.WhereFns, parseField2WhereFn(field)...)
		meta.OrderFns = append(meta.OrderFns, parseField2OrderFn(field)...)
	}

	for _, relation := range s.Relationships.HasOne {
		meta.PreloadFns = append(meta.PreloadFns, parseRelation2PreloadFn(relation)...)
	}
	for _, relation := range s.Relationships.HasMany {
		meta.PreloadFns = append(meta.PreloadFns, parseRelation2PreloadFn(relation)...)
	}

	return meta, nil
}

// parseField2WhereFn 字段解析为where
func parseField2WhereFn(field *schema.Field) []*method.Method {
	whereFns := make([]*method.Method, 0, 5)

	fieldType, err := fieldtype.GetConditionFieldBySchemaField(field)
	if err != nil {
		// TODO 后续补充完整字段之后需要返回error
		return whereFns
	}
	whereFns = append(whereFns, fieldType.GetWhereMethods()...)

	return whereFns
}

// parseField2OrderFn 字段解析为Order
func parseField2OrderFn(field *schema.Field) []*method.Method {
	orderFns := make([]*method.Method, 0, 5)
	// TODO
	if field.Name != "Sort" {
		return orderFns
	}
	fn := &method.Method{
		Field:      field,
		MethodName: "OrderBySort",
		Params:     []*method.Param{},
		Query:      fmt.Sprintf("%s ASC", field.DBName),
	}
	orderFns = append(orderFns, fn)
	return orderFns
}

// parseRelation2PreloadFn 解析Preload方法
func parseRelation2PreloadFn(relationship *schema.Relationship) []*method.Method {
	fns := make([]*method.Method, 0, 1)

	fn := &method.Method{
		MethodName: fmt.Sprintf("Preload%s", relationship.Name),
		Query:      relationship.Name,
	}
	fns = append(fns, fn)
	return fns
}
