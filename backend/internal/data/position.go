package data

import (
	"encoding/json"
	"errors"
	"strings"
)

var ErrInvalidPositionFormat = errors.New("invalid position format, expected: C|LW|RW|D|G")

type Position string

const (
	PositionC  Position = "C"
	PositionLW Position = "LW"
	PositionRW Position = "RW"
	PositionD  Position = "D"
	PositionG  Position = "G"
)

func (p *Position) UnmarshalJSON(jsonValue []byte) error {
	var s string
	if err := json.Unmarshal(jsonValue, &s); err != nil {
		return ErrInvalidPositionFormat
	}

	pos := Position(strings.ToUpper(s)) // force uppercase
	switch pos {
	case PositionC, PositionLW, PositionRW, PositionD, PositionG:
		*p = pos
		return nil
	}

	return ErrInvalidPositionFormat
}
