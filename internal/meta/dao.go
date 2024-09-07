package meta

type DaoStruct struct {
	QueryPkgName    string // 包名
	QueryStructList []*QueryStruct
}

func ParseDaoStruct(QueryPkgName string, queryStructList []*QueryStruct) *DaoStruct {
	return &DaoStruct{
		QueryPkgName:    QueryPkgName,
		QueryStructList: queryStructList,
	}
}
