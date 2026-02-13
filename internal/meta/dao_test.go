package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDaoStruct(t *testing.T) {
	queryStructs := []*QueryStruct{
		{
			Name: "User",
		},
		{
			Name: "Role",
		},
	}

	daoStruct := ParseDaoStruct("dao", queryStructs)

	assert.Equal(t, "dao", daoStruct.QueryPkgName)
	assert.Equal(t, queryStructs, daoStruct.QueryStructList)
}

