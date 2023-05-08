package utils

import (
	"errors"
	"strconv"
)

var ErrInvalidIntValue = errors.New("invalid int value")

func AsInt(v interface{}) (int, error) {
	switch v.(type) {
	case int:
		i, ok := v.(int)
		if !ok {
			return 0, ErrInvalidIntValue
		}

		return i, nil
	case float64:
		i, ok := v.(float64)
		if !ok {
			return 0, ErrInvalidIntValue
		}

		return int(i), nil
	case string:
		i, ok := v.(string)
		if !ok {
			return 0, ErrInvalidIntValue
		}

		return strconv.Atoi(i)
	}

	return 0, ErrInvalidIntValue
}
