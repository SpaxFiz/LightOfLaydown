package util

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/url"
	"reflect"
	"strings"
)

const DateLayout = "2006-01-02"

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

type TrimBound struct {
	Pre string
	Sur string
}

func UrlParamConstruct(target string, param map[string]interface{}) string {
	paramArr := make([]string, 0, len(param))
	for k, v := range param {
		current := k + "=" + url.QueryEscape(String(v))
		paramArr = append(paramArr, current)
	}
	if !strings.HasSuffix(target, "?") {
		target += "?"
	}
	return target + strings.Join(paramArr, "&")
}

func ParseBodyWithTrimAsJSON(body io.ReadCloser, trim *TrimBound, result interface{}) error {
	jsonStr, err := getBodySuffixTrim(body, trim)
	if err != nil {
		return err
	}
	if err := JSON.UnmarshalFromString(jsonStr, &result); err != nil {
		log.Error("parse body failed. err=%s", err.Error())
		return err
	}
	return nil
}

func getBodySuffixTrim(body io.ReadCloser, trim *TrimBound) (string, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	jsonStr := bytes.NewBuffer(b).String()
	if trim != nil {
		start := strings.Index(jsonStr, trim.Pre)
		end := indexLastRuneInString(jsonStr, []byte(trim.Sur)[0]) + 1
		jsonStr = jsonStr[start:end]
	}
	return jsonStr, nil
}

func indexLastRuneInString(s string, char byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == char {
			return i
		}
	}
	return -1
}

func RewriteStructJSON(input interface{}) (result interface{}, err error) {
	//defer func() {
	//	if e := recover(); e != nil {
	//		log.Errorf("RewriteStructJSON panic but recovered. err=%v", e)
	//		err = e.(error)
	//	}
	//}()

	rv := reflect.Indirect(reflect.ValueOf(input))
	switch rv.Kind() {
	case reflect.Struct:
		return extractLDTag(rv)
	case reflect.Slice:
		data := make([]interface{}, rv.Len())
		for j := 0; j < rv.Len(); j++ {
			if r, e := RewriteStructJSON(rv.Index(j).Interface()); e != nil {
				return nil, e
			} else {
				data[j] = r
			}
		}
		result = data
		return result, err
	}
	return nil, errors.New("unsupported type")
}

func extractLDTag(value reflect.Value) (map[string]interface{}, error) {
	elem := reflect.Indirect(value)
	if elem.Kind() != reflect.Struct {
		return nil, errors.New("input should like []interface{}")
	}
	r := make(map[string]interface{}, elem.NumField())
	for i := 0; i < elem.NumField(); i++ {
		key, ok := elem.Type().Field(i).Tag.Lookup("ld")
		if !ok {
			key = elem.Type().Field(i).Tag.Get("json")
		}
		if !elem.Field(i).IsValid() {
			r[key] = reflect.Zero(elem.Field(i).Type()).Interface()
		} else {
			r[key] = elem.Field(i).Interface()
		}
	}
	return r, nil
}
