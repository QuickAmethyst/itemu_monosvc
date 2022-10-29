package sql

import (
	"database/sql"
	"fmt"
)

var (
	ErrDestNil        = fmt.Errorf("dest cannot be nil")
	ErrWhereStructNil = fmt.Errorf("where struct cannot be nil")
	ErrNoRows         = sql.ErrNoRows
)
