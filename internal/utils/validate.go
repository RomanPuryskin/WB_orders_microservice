package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
)

func VaildateStructs[T any](someStruct T) error {
	validator := validator.New()
	err := validator.Struct(someStruct)
	if err != nil {
		return fmt.Errorf("[ValidateOrder|validate]: %w", err)
	}
	return nil
}

func ValidateUUID(id string) (uuid.UUID, error) {
	resUUID, err := uuid.FromString(id)
	if err != nil {
		return uuid.Nil, err
	}

	return resUUID, nil
}
