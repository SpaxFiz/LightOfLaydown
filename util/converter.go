package util

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"time"
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

type CustomDate string

func (c *CustomDate) UnmarshalJSON(b []byte) error {
	if val, err := strconv.ParseInt(string(b), 10, 64); err != nil {
		*c = "-"
	} else {
		*c = CustomDate(time.Unix(val/1000, 0).Format(DateLayout))
	}
	return nil
}

func AnyNum2Float64(in interface{}) (out float64, err error) {
	switch in.(type) {
	case float32:
		out = float64(in.(float32))
	case float64:
		out = in.(float64)
	case int:
		out = float64(in.(int))
	case int8:
		out = float64(in.(int8))
	case int16:
		out = float64(in.(int16))
	case int32:
		out = float64(in.(int32))
	case int64:
		out = float64(in.(int64))
	case uint8:
		out = float64(in.(uint8))
	case uint16:
		out = float64(in.(uint16))
	case uint32:
		out = float64(in.(uint32))
	case string:
		o, err := strconv.ParseFloat(in.(string), 10)
		if err != nil {
			return 0, errors.New("invalid input type, check if input is float")
		}
		out = o
	default:
		return 0, errors.New("invalid input type, check if input is float")
	}
	return out, nil
}
