package errs

import "fmt"

type ErrNotFound struct {
	entity string
	id     string
}

func NewErrNotFound(entityName string, id string) ErrNotFound {
	return ErrNotFound{entity: entityName, id: id}
}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("entity %s with id %s not found !", err.entity, err.id)
}
