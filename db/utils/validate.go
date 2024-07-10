package utils

import (
	"car-comparison-service/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func ValidateResultSuccess(result *gorm.DB) error {
	if result.Error != nil {
		return ValidateError(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.RECORD_NOT_FOUND
	}

	return nil
}

func ValidateError(err error) error {
	pgError, ok := err.(*pgconn.PgError)
	if ok {
		switch pgError.Code {
		case "23505", "23503":
			return errors.NewServiceError("CONSTRAINT_FAILURE: "+pgError.Message, 400)
		case "22P02":
			return errors.NewServiceError("INPUT_VALIDATION_FAILURE: "+pgError.Message, 400)
		case "40001":
			return errors.TRANSACTION_SERIALIZATION_ERROR
		}
	}

	if err.Error() == "record not found" {
		return errors.RECORD_NOT_FOUND
	}

	return errors.NewServiceError(err.Error(), 500)
}
