package util

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

type TolerantFloat float64

func (t *TolerantFloat) UnmarshalJSON(b []byte) error {
	if val, err := strconv.ParseFloat(string(b), 64); err != nil {
		*t = 0
	} else {
		*t = TolerantFloat(val)
	}
	return nil
}

func String(value interface{}) string {
	switch value := value.(type) {
	case int:
		return strconv.Itoa(value)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case json.Number:
		return value.String()
	case []byte:
		return string(value)
	case []interface{}, []int64, []float64, []string:
		v, _ := json.Marshal(value)
		return string(v)
	case nil:
		return ""
	}
	return fmt.Sprintf("%v", value)
}

func Float64FromBytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
