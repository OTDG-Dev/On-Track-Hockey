package data

import (
	"encoding/json"
	"errors"
)

var ErrInvalidShootCatches = errors.New("invalid position format, expected: L|R")

type ShootsCatches string

const (
	ShootsCatchesL ShootsCatches = "L"
	ShootsCatchesR ShootsCatches = "R"
)

func (sc *ShootsCatches) UnmarshalJSON(jsonValue []byte) error {
	var s string
	if err := json.Unmarshal(jsonValue, &s); err != nil {
		return ErrInvalidPositionFormat
	}

	shct := ShootsCatches(s)
	switch shct {
	case ShootsCatchesL, ShootsCatchesR:
		*sc = shct
		return nil
	}

	return ErrInvalidShootCatches
}
