package auth

type Permission = int8

const (
	READ Permission = iota
	WRITE
	DELETE
)
