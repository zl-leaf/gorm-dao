package tpl

const (
	QueryStructTpl = `
package {{.QueryPkgName}}

type {{.StructName}}Dao struct {
	*gorm.DB
}

func new{{.Name}}Dao(db *gorm.DB) *{{.StructName}}Dao {
	return &{{.StructName}}Dao{db}
}

func (d *{{.StructName}}Dao) WithContext(ctx context.Context) *{{.StructName}}Dao {
	return &{{.StructName}}Dao{d.DB.WithContext(ctx)}
}

` + QueryMethodTpl + CRUDTpl

	QueryMethodTpl = `
{{range .WhereFns -}}
func (d *{{$.StructName}}Dao) {{.MethodName}}({{.Params.InputArgs}}) *{{$.StructName}}Dao {
	d.DB = d.DB.Where({{.Query}})
	return d
}
{{end}}

{{range .OrderFns -}}
func (d *{{$.StructName}}Dao) {{.MethodName}}() *{{$.StructName}}Dao {
	d.DB = d.DB.Order("{{.Query}}")
	return d
}
{{end}}

func(d *{{.StructName}}Dao) Limit(limit int) *{{.StructName}}Dao {
	d.DB = d.DB.Limit(limit)
	return d
}


func (d *{{.StructName}}Dao) Offset(offset int) *{{.StructName}}Dao {
	d.DB = d.DB.Offset(offset)
	return d
}

func(d *{{.StructName}}Dao) Page(page, pageSize int) *{{.StructName}}Dao {
	if page <= 0 {
		page = 1
	}
	if pageSize <=0 {
		pageSize = 10
	}
	d.DB = d.DB.Offset((page-1) * pageSize)
	d.DB = d.DB.Limit(pageSize)
	return d
}
`

	CRUDTpl = `
func (d *{{.StructName}}Dao) GetByID(id int64) (*model.{{.Name}}, error) {
	return d.First(id)
}

func (d *{{.StructName}}Dao) Find() ([]*model.{{.Name}}, int64, error) {
	var list []*model.{{.Name}}
	var total int64
	var query = d.Model(model.{{.Name}}{})

	if err := query.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (d *{{.StructName}}Dao) First(conds ...interface{}) (*model.{{.Name}}, error) {
	var record = &model.{{.Name}}{}
	result := d.DB.First(record, conds...)
	if result.Error != nil {
		return nil, result.Error
	}
	return record, nil
}

func (d *{{.StructName}}Dao) Create(record *model.{{.Name}}) error {
	return d.DB.Create(record).Error
}

func (d *{{.StructName}}Dao) Save(record *model.{{.Name}}) error {
	return d.DB.Save(record).Error
}

func (d *{{.StructName}}Dao) Delete(conds ...interface{}) error {
	return d.DB.Delete(&model.{{.Name}}{}, conds...).Error
}
`
)
