package qbit

type TypeMapper interface {
	Convert(colType string, tagType string) *Type
}

type MysqlTypeMapper struct {
}

func (t MysqlTypeMapper) Convert(colType string, tagType string) *Type {
	return nil
}

type SqliteTypeMapper struct {
}

func (t SqliteTypeMapper) Convert(colType string, tagType string) *Type {
	return nil
}

type PostgresTypeMapper struct {
}

func (t PostgresTypeMapper) Convert(colType string, tagType string) *Type {
	return nil
}
