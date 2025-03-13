package domainc

import "database/sql"

type Query interface {
	RunQuery(query string, values ...interface{}) (sql.Result, error)
	GetData(query string, values ...interface{}) (sql.Result, error)
}
