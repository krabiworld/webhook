package codec

import (
	"fmt"
	"strconv"
)

type BoolCodec struct{}

func (c BoolCodec) Encode(v bool) (string, error) {
	return strconv.FormatBool(v), nil
}

func (c BoolCodec) Decode(s string) (bool, error) {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, fmt.Errorf("BoolCodec: failed to parse %q as bool: %w", s, err)
	}
	return v, nil
}
