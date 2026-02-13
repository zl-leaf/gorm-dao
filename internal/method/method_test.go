package method

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParamsInputArgs(t *testing.T) {
	params := Params{
		{
			Name:    "user",
			Type:    "User",
			Package: "model",
		},
		{
			Name:      "userPtr",
			Type:      "User",
			Package:   "model",
			IsPointer: true,
		},
		{
			Name:    "users",
			Type:    "User",
			Package: "model",
			IsArray: true,
		},
	}

	input := params.InputArgs()
	assert.Equal(t, "model.user User,*model.userPtr User,model.users []User", input)
}

func TestParamsQueryArgs(t *testing.T) {
	params := Params{
		{Name: "id"},
		{Name: "name"},
	}

	queryArgs := params.QueryArgs()
	assert.Equal(t, "id,name", queryArgs)
}

func TestToArgsName(t *testing.T) {
	tests := []struct {
		fieldName string
		want      string
	}{
		{"ID", "ID"},
		{"Name", "name"},
		{"UserName", "userName"},
		{"X", "x"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, ToArgsName(tt.fieldName))
	}
}

