package gormdao

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/zl-leaf/gorm-dao/internal/meta"
	"github.com/zl-leaf/gorm-dao/internal/tpl"
	"golang.org/x/tools/imports"
	"gorm.io/gorm/schema"
)

type Generator struct {
	outputPath string
	models     []interface{}
}

func NewGenerator(outputPath string) *Generator {
	return &Generator{
		outputPath: outputPath,
	}
}

func (g *Generator) Apply(models ...interface{}) {
	g.models = models
}

func (g *Generator) Generate() error {
	// 处理dao目录，生成文件夹
	var err error
	g.outputPath, err = filepath.Abs(g.outputPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(g.outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	modelSchemaList := make([]*schema.Schema, len(g.models))
	for i, m := range g.models {
		modelSchema, err := schema.Parse(&m, &sync.Map{}, schema.NamingStrategy{})
		if err != nil {
			return err
		}
		modelSchemaList[i] = modelSchema
	}

	queryStructList, err := g.generateQueryStructs(modelSchemaList)
	if err != nil {
		return err
	}
	err = g.generateDaoStruct(queryStructList)
	if err != nil {
		return err
	}

	return nil
}

// pkgName 获取包名
func (g *Generator) pkgName() string {
	return filepath.Base(g.outputPath)
}

// generateDaoStruct 生成dao主文件
func (g *Generator) generateDaoStruct(queryStructList []*meta.QueryStruct) error {
	structTpl, err := template.New("dao").Parse(tpl.DaoStructTpl)
	if err != nil {
		return err
	}

	daoStruct := meta.ParseDaoStruct(g.pkgName(), queryStructList)
	// 生成文件
	var buf bytes.Buffer
	err = structTpl.Execute(&buf, daoStruct)
	if err != nil {
		return err
	}
	fileName := filepath.Join(g.outputPath, "dao.go")
	err = g.output(fileName, buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// generateQueryStructs 生成dao查询文件
func (g *Generator) generateQueryStructs(modelSchemaList []*schema.Schema) ([]*meta.QueryStruct, error) {
	structTpl, err := template.New("query").Parse(tpl.QueryStructTpl)
	if err != nil {
		return nil, err
	}

	queryStructList := make([]*meta.QueryStruct, len(modelSchemaList))
	for i, modelSchema := range modelSchemaList {
		queryStruct, err := meta.ParseQueryStruct(g.pkgName(), modelSchema)
		if err != nil {
			return nil, err
		}
		queryStructList[i] = queryStruct

		// 生成文件
		var buf bytes.Buffer
		err = structTpl.Execute(&buf, queryStruct)
		if err != nil {
			return nil, err
		}
		fileName := g.outputPath + "/" + queryStruct.Schema.Table + ".go"
		err = g.output(fileName, buf.Bytes())
		if err != nil {
			return nil, err
		}
	}
	return queryStructList, nil
}

// output 输出go文件
func (g *Generator) output(fileName string, content []byte) error {
	result, err := imports.Process(fileName, content, nil)
	if err != nil {
		lines := strings.Split(string(content), "\n")
		errLine, _ := strconv.Atoi(strings.Split(err.Error(), ":")[1])
		startLine, endLine := errLine-5, errLine+5
		fmt.Println("Format fail:", errLine, err)
		if startLine < 0 {
			startLine = 0
		}
		if endLine > len(lines)-1 {
			endLine = len(lines) - 1
		}
		for i := startLine; i <= endLine; i++ {
			fmt.Println(i, lines[i])
		}
		return fmt.Errorf("cannot format file: %w", err)
	}
	return os.WriteFile(fileName, result, 0640)
}
