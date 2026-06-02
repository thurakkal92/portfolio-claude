package contact

import "github.com/jackc/pgx/v5"

// pgxNoRows returns the pgx sentinel; isolated to avoid importing pgx in service.go's main flow.
func pgxNoRows() error { return pgx.ErrNoRows }
