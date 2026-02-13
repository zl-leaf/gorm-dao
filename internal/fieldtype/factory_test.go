package fieldtype

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestGetConditionFieldBySchemaField_UnsupportedType(t *testing.T) {
	s := parseDemoSchema(t)

	field := s.FieldsByName["Unsupported"]
	condField, err := GetConditionFieldBySchemaField(field)
	assert.Error(t, err)
	assert.Nil(t, condField)
}

