package errs

import "fmt"

type ErrDuplicate struct {
	entity string
	field  string
}

func NewErrDuplicate(entityName string, field string) ErrDuplicate {
	return ErrDuplicate{entity: entityName, field: field}
}

func (err ErrDuplicate) Error() string {
	if err.field != "" {
		return "duplicate entry"

	} else {
		return fmt.Sprintf("field %s on %s is duplicate", err.field, err.entity)
	}
}
