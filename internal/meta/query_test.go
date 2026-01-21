package meta

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DemoModel struct {
	gorm.Model

	HasOneModel   HasOneModel
	HasManyModels []HasManyModel
}

type HasOneModel struct {
	gorm.Model
	DemoModelID uint
}

type HasManyModel struct {
	gorm.Model
	DemoModelID uint
}

func TestParseQueryStruct(t *testing.T) {
	demoModel := &DemoModel{}
	modelSchema, err := schema.Parse(&demoModel, &sync.Map{}, schema.NamingStrategy{})
	assert.Nil(t, err)
	queryStruct, err := ParseQueryStruct("demo", modelSchema)
	assert.Nil(t, err)

	assert.NotNil(t, queryStruct)
	assert.NotEmpty(t, queryStruct.WhereFns)
	assert.NotEmpty(t, queryStruct.PreloadFns)
	assert.Equal(t, "PreloadHasOneModel", queryStruct.PreloadFns[0].MethodName)
	assert.Equal(t, "PreloadHasManyModels", queryStruct.PreloadFns[1].MethodName)
}
