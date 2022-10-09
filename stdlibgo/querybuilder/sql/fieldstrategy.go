package sql

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

var (
	// https://github.com/golang/lint/blob/master/lint.go#L770
	commonInitialisms         = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	commonInitialismsReplacer *strings.Replacer
	statementSuffixes         = []string{"_not_eq", "_like", "_gte", "_lte", "_not_in", "_in", "_is_null"}
)

func init() {
	commonInitialismsForReplacer := make([]string, 0, len(commonInitialisms))
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(
			commonInitialismsForReplacer,
			initialism,
			cases.Title(language.Und, cases.NoLower).String(initialism),
		)
	}

	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}

type FieldStrategy string

func (f FieldStrategy) IsEqualStatement() bool {
	return !f.IsNotEqualStatement() &&
		!f.IsLikeStatement() &&
		!f.IsGreaterThanEqualStatement() &&
		!f.IsLessThanEqualStatement() &&
		!f.IsInStatement() &&
		!f.IsNotInStatement() &&
		!f.IsNull()
}

func (f FieldStrategy) IsNotEqualStatement() bool {
	return strings.HasSuffix(string(f), "NotEQ")
}

func (f FieldStrategy) IsLikeStatement() bool {
	return strings.HasSuffix(string(f), "Like")
}

func (f FieldStrategy) IsGreaterThanEqualStatement() bool {
	return strings.HasSuffix(string(f), "GTE")
}

func (f FieldStrategy) IsLessThanEqualStatement() bool {
	return strings.HasSuffix(string(f), "LTE")
}

func (f FieldStrategy) IsInStatement() bool {
	return strings.HasSuffix(string(f), "IN") && !f.IsNotInStatement()
}

func (f FieldStrategy) IsNotInStatement() bool {
	return strings.HasSuffix(string(f), "NotIN")
}

func (f FieldStrategy) IsNull() bool {
	return strings.HasSuffix(string(f), "IsNULL")
}

func (f FieldStrategy) ColumnName() string {
	return ColumnName(string(f))
}
