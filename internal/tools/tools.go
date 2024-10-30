package tools

import (
	"bytes"
	"encoding/binary"
)

func DataConvert[T int | uint8 | uint16 | int8 | int16 | int64 | int32 | float32 | float64 | bool | string](data []byte) (T, error) {
	var temp T
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &temp)
	return temp, err
}
