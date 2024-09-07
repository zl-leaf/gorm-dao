package tpl

const (
	DaoStructTpl = `
package {{.QueryPkgName}}

type Dao struct {
	db        *gorm.DB
}


func NewDao(db *gorm.DB) *Dao {
	return &Dao{
		db:        db,
	}
}

{{range .QueryStructList -}}
	func (d *Dao) {{.Name}}Dao (ctx context.Context) *{{.StructName}}Dao {
		return new{{.Name}}Dao(d.db).WithContext(ctx)
	}
{{end}}

func (d *Dao) Transaction(ctx context.Context, fc func(tx *Dao) error, opts ...*sql.TxOptions) (err error) {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newDao := NewDao(tx)
		return fc(newDao)
	}, opts...)
}
`
)
