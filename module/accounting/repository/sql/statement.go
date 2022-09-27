package sql

type AccountClassStatement struct {
	ID int64
}

type AccountGroupStatement struct {
	ID             int64
	ParentID       int64
	ParentIDIsNULL bool
}
