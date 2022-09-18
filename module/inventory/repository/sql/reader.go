package sql

import (
	"context"
	"fmt"
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	qb "github.com/QuickAmethyst/monosvc/stdlibgo/querybuilder/sql"
	"github.com/QuickAmethyst/monosvc/stdlibgo/sql"

	"github.com/QuickAmethyst/monosvc/module/inventory/domain"
)

type Reader interface {
	GetUomList(ctx context.Context, stmt UomStatement, p qb.Paging) (result []domain.Uom, paging qb.Paging, err error)
}

type reader struct {
	db sql.DB
}

func (r *reader) GetUomList(ctx context.Context, stmt UomStatement, p qb.Paging) (result []domain.Uom, paging qb.Paging, err error) {
	result = make([]domain.Uom, 0)
	paging = p
	paging.Normalize()

	fromClause := "FROM uoms"
	limitClause, limitClauseArgs := p.BuildQuery()
	whereClause, whereClauseArgs, err := qb.NewWhereClause(stmt)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetUomListFailed, "Failed on build where clause")
		return
	}

	selectQuery := fmt.Sprintf("SELECT id, name, description, decimal %s %s %s", fromClause, whereClause, limitClause)
	countQuery := fmt.Sprintf("SELECT COUNT(*) %s %s", fromClause, whereClause)

	if err = r.db.SelectContext(ctx, &result, r.db.Rebind(selectQuery), append(whereClauseArgs, limitClauseArgs...)...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetUomListFailed, "Failed on select uom")
		return
	}

	if err = r.db.GetContext(ctx, &paging.Total, r.db.Rebind(countQuery), whereClauseArgs...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetUomListCountFailed, "Failed on select count uom")
		return
	}

	return
}

func NewReader(opt *Options) Reader {
	return &reader{opt.SlaveDB}
}
