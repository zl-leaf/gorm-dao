package fieldtype

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zl-leaf/gorm-dao/internal/method"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type demoNumberStringModel struct {
	gorm.Model

	Name        string
	Age         int
	Score       float64
	Unsupported bool
}

type demoTimeModel struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 指针类型字段，用于测试 IsNull / IsNotNull
type demoPointerModel struct {
	ID        uint
	Name      *string
	Age       *int
	DeletedAt *time.Time
}

func parseDemoSchema(t *testing.T) *schema.Schema {
	t.Helper()

	m := &demoNumberStringModel{}
	s, err := schema.Parse(m, &sync.Map{}, schema.NamingStrategy{})
	assert.NoError(t, err)
	return s
}

func TestGetConditionFieldBySchemaField_StringField(t *testing.T) {
	s := parseDemoSchema(t)

	nameField := s.FieldsByName["Name"]
	condField, err := GetConditionFieldBySchemaField(nameField)
	assert.NoError(t, err)

	sf, ok := condField.(*StringField)
	assert.True(t, ok, "expected StringField")

	methods := sf.GetWhereMethods()
	// string 类型应该包含等于、不等于、like、前缀 like、not like、in、not in 共 7 个方法
	assert.Len(t, methods, 7)

	methodNames := make([]string, 0, len(methods))
	for _, m := range methods {
		methodNames = append(methodNames, m.MethodName)
	}

	assert.Contains(t, methodNames, "WhereNameEq")
	assert.Contains(t, methodNames, "WhereNameLike")
	assert.Contains(t, methodNames, "WhereNamePrefixLike")
	assert.Contains(t, methodNames, "WhereNameNotLike")
	assert.Contains(t, methodNames, "WhereNameIn")
	assert.Contains(t, methodNames, "WhereNameNotIn")
}

func TestGetConditionFieldBySchemaField_NumberField(t *testing.T) {
	s := parseDemoSchema(t)

	ageField := s.FieldsByName["Age"]
	condField, err := GetConditionFieldBySchemaField(ageField)
	assert.NoError(t, err)

	nf, ok := condField.(*NumberField)
	assert.True(t, ok, "expected NumberField")

	methods := nf.GetWhereMethods()
	// number 类型包含等于、不等于、in、not in、lt、lte、gt、gte、between 共 9 个方法
	assert.Len(t, methods, 9)

	methodNames := make([]string, 0, len(methods))
	for _, m := range methods {
		methodNames = append(methodNames, m.MethodName)
	}

	assert.Contains(t, methodNames, "WhereAgeEq")
	assert.Contains(t, methodNames, "WhereAgeLt")
	assert.Contains(t, methodNames, "WhereAgeLte")
	assert.Contains(t, methodNames, "WhereAgeGt")
	assert.Contains(t, methodNames, "WhereAgeGte")
	assert.Contains(t, methodNames, "WhereAgeBetween")
	assert.Contains(t, methodNames, "WhereAgeIn")
	assert.Contains(t, methodNames, "WhereAgeNotIn")
}

func TestGetConditionFieldBySchemaField_TimeField(t *testing.T) {
	model := &demoTimeModel{}
	s, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	assert.NoError(t, err)

	createdAtField := s.FieldsByName["CreatedAt"]
	condField, err := GetConditionFieldBySchemaField(createdAtField)
	assert.NoError(t, err)

	tf, ok := condField.(*TimeField)
	assert.True(t, ok, "expected TimeField")

	methods := tf.GetWhereMethods()
	// time 类型包含 eq、gt、gte、lt、lte、between 共 6 个方法
	assert.Len(t, methods, 6)

	methodNames := make([]string, 0, len(methods))
	for _, fn := range methods {
		methodNames = append(methodNames, fn.MethodName)
	}

	assert.Contains(t, methodNames, "WhereCreatedAtEq")
	assert.Contains(t, methodNames, "WhereCreatedAtGt")
	assert.Contains(t, methodNames, "WhereCreatedAtGte")
	assert.Contains(t, methodNames, "WhereCreatedAtLt")
	assert.Contains(t, methodNames, "WhereCreatedAtLte")
	assert.Contains(t, methodNames, "WhereCreatedAtBetween")
}

func TestGetConditionFieldBySchemaField_PointerStringAddsIsNullIsNotNull(t *testing.T) {
	model := &demoPointerModel{}
	s, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	assert.NoError(t, err)

	nameField := s.FieldsByName["Name"]
	condField, err := GetConditionFieldBySchemaField(nameField)
	assert.NoError(t, err)

	sf, ok := condField.(*StringField)
	assert.True(t, ok, "expected StringField")

	methods := sf.GetWhereMethods()
	// string 7 个 + IsNull + IsNotNull = 9
	assert.Len(t, methods, 9)

	methodNames := make([]string, 0, len(methods))
	for _, fn := range methods {
		methodNames = append(methodNames, fn.MethodName)
	}
	assert.Contains(t, methodNames, "WhereNameIsNull")
	assert.Contains(t, methodNames, "WhereNameIsNotNull")
}

func TestGetConditionFieldBySchemaField_PointerNumberAddsIsNullIsNotNull(t *testing.T) {
	model := &demoPointerModel{}
	s, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	assert.NoError(t, err)

	ageField := s.FieldsByName["Age"]
	condField, err := GetConditionFieldBySchemaField(ageField)
	assert.NoError(t, err)

	nf, ok := condField.(*NumberField)
	assert.True(t, ok, "expected NumberField")

	methods := nf.GetWhereMethods()
	// number 9 个 + IsNull + IsNotNull = 11
	assert.Len(t, methods, 11)

	methodNames := make([]string, 0, len(methods))
	for _, fn := range methods {
		methodNames = append(methodNames, fn.MethodName)
	}
	assert.Contains(t, methodNames, "WhereAgeIsNull")
	assert.Contains(t, methodNames, "WhereAgeIsNotNull")
}

// 指针类型字段的 eq 等方法，Params 的 Type 应为元素类型（如 int），而非 *int
func TestGetConditionFieldBySchemaField_PointerFieldParamsUseElementType(t *testing.T) {
	model := &demoPointerModel{}
	s, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	assert.NoError(t, err)

	ageField := s.FieldsByName["Age"]
	condField, err := GetConditionFieldBySchemaField(ageField)
	assert.NoError(t, err)

	nf, ok := condField.(*NumberField)
	assert.True(t, ok, "expected NumberField")

	methods := nf.GetWhereMethods()
	var eqMethod *method.Method
	for _, m := range methods {
		if m.MethodName == "WhereAgeEq" {
			eqMethod = m
			break
		}
	}
	assert.NotNil(t, eqMethod, "WhereAgeEq method should exist")
	assert.Len(t, eqMethod.Params, 1, "WhereAgeEq should have one param")
	assert.Equal(t, "int", eqMethod.Params[0].Type, "pointer *int field should have param Type int, not *int")
}

func TestGetConditionFieldBySchemaField_PointerTimeAddsIsNullIsNotNull(t *testing.T) {
	model := &demoPointerModel{}
	s, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	assert.NoError(t, err)

	field := s.FieldsByName["DeletedAt"]
	condField, err := GetConditionFieldBySchemaField(field)
	assert.NoError(t, err)

	tf, ok := condField.(*TimeField)
	assert.True(t, ok, "expected TimeField")

	methods := tf.GetWhereMethods()
	// time 6 个 + IsNull + IsNotNull = 8
	assert.Len(t, methods, 8)

	methodNames := make([]string, 0, len(methods))
	for _, fn := range methods {
		methodNames = append(methodNames, fn.MethodName)
	}
	assert.Contains(t, methodNames, "WhereDeletedAtIsNull")
	assert.Contains(t, methodNames, "WhereDeletedAtIsNotNull")
}

func TestGetConditionFieldBySchemaField_UnsupportedType(t *testing.T) {
	s := parseDemoSchema(t)

	field := s.FieldsByName["Unsupported"]
	condField, err := GetConditionFieldBySchemaField(field)
	assert.Error(t, err)
	assert.Nil(t, condField)
}
