package model

import "fmt"

// TODO this is more sync related

const (
	ChangeSourceUnspecified = 1
	ChangeSourceSync        = 2
	ChangeSourceDecryption  = 3

	ChangeTypeCreate = 1
	ChangeTypeUpdate = 2
	ChangeTypeDelete = 3
)

type ChangeType int

func ChangeTypeFromInt(input int) (ChangeType, error) {
	switch input {
	case ChangeTypeCreate, ChangeTypeUpdate, ChangeTypeDelete:
		return ChangeType(input), nil
	default:
		return 0, fmt.Errorf("unknown change type %d", input)
	}
}

type ChangeSourceType int

func ChangeSourceTypeFromInt(input int) (ChangeSourceType, error) {
	switch input {
	case ChangeSourceUnspecified, ChangeSourceSync, ChangeSourceDecryption:
		return ChangeSourceType(input), nil
	default:
		return ChangeSourceUnspecified, fmt.Errorf("unknown change source type %d", input)
	}
}

type ItemChange struct {
	Item
	ItemID       string
	ChangeSource ChangeSourceType
	ChangeType   ChangeType
}
