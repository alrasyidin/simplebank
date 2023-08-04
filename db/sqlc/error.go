package db

import (
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
)

var ErrRecordNotFound = pgx.ErrNoRows
var ErrUniqueViolations = pgerrcode.UniqueViolation
var ErrForeignKeyViolations = pgerrcode.ForeignKeyViolation
