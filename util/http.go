package util

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
)

const ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36"

func EasyGet(target string, reqParam map[string]interface{}, trim *TrimBound, receiver interface{}, header map[string]string) error {
	url := UrlParamConstruct(target, reqParam)
	fmt.Println(url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	if header == nil {
		req.Header.Set("user-agent", ua)
		req.Header.Add("Accept-Charset", "utf-8")
	} else {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error occur when request  get. err=%s", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("response http code is not 200. got %d", resp.StatusCode)
	}

	receiverType := reflect.TypeOf(receiver)
	if receiverType.Kind() == reflect.Ptr && receiverType.Elem().Kind() == reflect.String {
		*(receiver.(*string)), err = getBodySuffixTrim(resp.Body, trim)
		if err != nil {
			return err
		}
		return nil
	}
	if err := ParseBodyWithTrimAsJSON(resp.Body, trim, &receiver); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func EasyPost(target string, payload map[string]interface{}, header map[string]string, receiver interface{}) error {
	client := &http.Client{}
	body, err := JSON.Marshal(payload)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", target, bytes.NewBuffer(body))
	if header == nil {
		req.Header.Set("user-agent", ua)
		req.Header.Add("Accept-Charset", "utf-8")
		req.Header.Set("Content-Type", "application/json")
	} else {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	resp, err2 := client.Do(req)
	if err2 != nil {
		log.Errorf("error occur when request post. err=%s", err.Error())
		return err2
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("response http code is not 200. got %d", resp.StatusCode)
	}
	raw, _ := ioutil.ReadAll(resp.Body)

	if err := JSON.Unmarshal(raw, &receiver); err != nil {
		return err
	}
	return nil
}
