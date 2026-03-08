package data

import (
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// https://www.postgresql.org/docs/8.2/errcodes-appendix.html
const (
	RestrictViolationFKeyErr = "23001"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrFKeyViolation  = errors.New("foreign key violation")
)

// ExecDeleteErrors converts database delete errors into
// application friendly error messages.
//
//	fromTable is the table you are deleting from
func ExecDeleteErrors(err error, fromTable string) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case RestrictViolationFKeyErr:
			return fmt.Errorf("%w: cannot delete %s it contains one or more %s",
				ErrFKeyViolation,
				fromTable,
				pqErr.Table,
			)
		}
	}
	return err
}
