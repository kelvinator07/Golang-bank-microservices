package db

import (
	"database/sql"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrRecordNotFound = sql.ErrNoRows

var ErrUniqueViolation = sql.ErrConnDone
