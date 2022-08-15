package sql

func NewWhereClause(stmt interface{}) (whereClause string, args []interface{}, err error) {
	condition := Condition{stmt}
	whereClause, args, err = condition.BuildQuery()
	return
}
