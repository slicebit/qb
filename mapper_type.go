package qbit

type TypeMapper interface {
	Convert() *Type
}

type MysqlTypeMapper struct {
}

type SqliteTypeMapper struct {
}

type PostgresTypeMapper struct {
}
