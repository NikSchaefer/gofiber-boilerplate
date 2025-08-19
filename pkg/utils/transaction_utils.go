package utils

import (
	"errors"

	"github.com/NikSchaefer/go-fiber/ent"
)

// RollbackTx rolls back the transaction and returns the error
func RollbackTx(tx *ent.Tx, err error) error {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}
	return err
}
