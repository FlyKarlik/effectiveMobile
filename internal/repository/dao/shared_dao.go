package dao

import (
	"database/sql"

	"github.com/FlyKarlik/effectiveMobile/internal/domain"
	"github.com/FlyKarlik/effectiveMobile/pkg/database/postgres"
)

type PaginationDAO struct {
	Limit  sql.NullInt64
	Offset sql.NullInt64
}

func (p *PaginationDAO) FromDomain(domain domain.Pagination) {
	if domain.Limit != 0 {
		p.Limit = postgres.ToNullInt64(&domain.Limit)
	}

	if domain.Offset != 0 {
		p.Offset = postgres.ToNullInt64(&domain.Offset)
	}
}
