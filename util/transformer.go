package util

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/url"
	"reflect"
	"strings"
)

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

func RewriteDataStructJSON(input interface{}) (result interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Errorf("RewriteDataStructJSON panic but recovered. err=%v", e)
			err = e.(error)
		}
	}()
	rv := reflect.ValueOf(input)
	data := make([]map[string]interface{}, rv.Len())
	for j := 0; j < rv.Len(); j++ {
		rt := reflect.TypeOf(rv.Index(j).Interface()).Elem()
		innerRv := reflect.ValueOf(rv.Index(j).Interface()).Elem()
		result := make(map[string]interface{}, rt.NumField())
		for i := 0; i < rt.NumField(); i++ {
			key := rt.Field(i).Tag.Get("ld")
			if !innerRv.Field(i).IsValid() {
				result[key] = reflect.Zero(rt.Field(i).Type).Interface()
			} else {
				result[key] = innerRv.Field(i).Interface()
			}
		}
		data[j] = result
	}
	result = err
	return result, err
}
